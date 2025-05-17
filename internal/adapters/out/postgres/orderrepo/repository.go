package orderrepo

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/uow"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.OrderRepository = &Repository{}

type Repository struct {
	tracker uow.Tracker
}

func NewRepository(tracker uow.Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}
	return &Repository{tracker: tracker}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *order.Order) error {
	return r.saveWithTx(ctx, DomainToDTO(aggregate), true)
}

func (r *Repository) Update(ctx context.Context, aggregate *order.Order) error {
	return r.saveWithTx(ctx, DomainToDTO(aggregate), false)
}

func (r *Repository) saveWithTx(ctx context.Context, dto OrderDTO, isCreate bool) error {
	r.tracker.Track(DtoToDomain(dto))

	inTx := r.tracker.InTx()
	if !inTx {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	session := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true})
	var err error
	if isCreate {
		err = session.Create(&dto).Error
	} else {
		err = session.Save(&dto).Error
	}
	if err != nil {
		return err
	}

	if !inTx {
		if err := r.tracker.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*order.Order, error) {
	var dto OrderDTO
	tx := r.txOrDb()

	if err := tx.WithContext(ctx).Preload(clause.Associations).First(&dto, ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return DtoToDomain(dto), nil
}

func (r *Repository) GetFirstInCreatedStatus(ctx context.Context) (*order.Order, error) {
	var dto OrderDTO
	err := r.txOrDb().WithContext(ctx).
		Preload(clause.Associations).
		Where("status = ?", order.StatusCreated).
		First(&dto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewObjectNotFoundError("Order (Status=Created)", nil)
		}
		return nil, err
	}
	return DtoToDomain(dto), nil
}

func (r *Repository) GetAllInAssignedStatus(ctx context.Context) ([]*order.Order, error) {
	var dtos []OrderDTO
	err := r.txOrDb().WithContext(ctx).
		Preload(clause.Associations).
		Where("status = ?", order.StatusAssigned).
		Find(&dtos).Error

	if err != nil {
		return nil, err
	}
	if len(dtos) == 0 {
		return nil, errs.NewObjectNotFoundError("Assigned orders", nil)
	}

	return mapDtosToAggregates(dtos), nil
}

func (r *Repository) txOrDb() *gorm.DB {
	if tx := r.tracker.Tx(); tx != nil {
		return tx
	}
	return r.tracker.Db()
}

func mapDtosToAggregates(dtos []OrderDTO) []*order.Order {
	aggregates := make([]*order.Order, len(dtos))
	for i := range dtos {
		aggregates[i] = DtoToDomain(dtos[i])
	}
	return aggregates
}
