package postgres

import (
	"context"
	"delivery/internal/adapters/out/postgres/courierrepo"
	"delivery/internal/adapters/out/postgres/orderrepo"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox"
	"errors"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

var _ ports.UnitOfWork = &UnitOfWork{}

type UnitOfWork struct {
	tx                *gorm.DB
	db                *gorm.DB
	trackedAggregates []ddd.AggregateRoot
	courierRepository ports.CourierRepository
	orderRepository   ports.OrderRepository
}

func (u *UnitOfWork) Rollback() error {
	if u.tx != nil {
		err := u.tx.Rollback().Error
		u.tx = nil
		return err
	}
	return nil
}

func NewUnitOfWork(db *gorm.DB) (ports.UnitOfWork, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	uow := &UnitOfWork{
		db: db,
	}

	courierRepo, err := courierrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.courierRepository = courierRepo

	orderRepo, err := orderrepo.NewRepository(uow)
	if err != nil {
		return nil, err
	}
	uow.orderRepository = orderRepo

	return uow, nil
}

func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}

func (u *UnitOfWork) Db() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) InTx() bool {
	return u.tx != nil
}

func (u *UnitOfWork) Track(agg ddd.AggregateRoot) {
	u.trackedAggregates = append(u.trackedAggregates, agg)
}

func (u *UnitOfWork) CourierRepository() ports.CourierRepository {
	return u.courierRepository
}

func (u *UnitOfWork) OrderRepository() ports.OrderRepository {
	return u.orderRepository
}

func (u *UnitOfWork) Begin(ctx context.Context) {
	u.tx = u.db.WithContext(ctx).Begin()
}

func (u *UnitOfWork) Commit(ctx context.Context) error {
	if u.tx == nil {
		return errs.NewValueIsRequiredError("cannot commit without transaction")
	}

	committed := false
	defer func() {
		if !committed {
			if err := u.tx.WithContext(ctx).Rollback().Error; err != nil && !errors.Is(err, gorm.ErrInvalidTransaction) {
				log.Error(err)
			}
			u.clearTx()
		}
	}()

	if err := u.persistDomainEvents(ctx, u.tx); err != nil {
		return err
	}

	if err := u.tx.WithContext(ctx).Commit().Error; err != nil {
		return err
	}
	committed = true
	u.clearTx()

	return nil
}

func (u *UnitOfWork) clearTx() {
	u.tx = nil
	u.trackedAggregates = nil
}

func (u *UnitOfWork) persistDomainEvents(ctx context.Context, tx *gorm.DB) error {
	for _, agg := range u.trackedAggregates {
		outboxEvents, err := outbox.EncodeDomainEvents(agg.GetDomainEvents())
		if err != nil {
			return err
		}
		if len(outboxEvents) > 0 {
			if err := tx.WithContext(ctx).Create(&outboxEvents).Error; err != nil {
				return err
			}
		}
		agg.ClearDomainEvents()
	}
	return nil
}
