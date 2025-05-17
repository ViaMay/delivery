package courierrepo

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/uow"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CourierRepository = &Repository{}

type Repository struct {
	tracker uow.Tracker
}

func NewRepository(tracker uow.Tracker) (*Repository, error) {
	if tracker == nil {
		return nil, errs.NewValueIsRequiredError("tracker")
	}
	return &Repository{tracker: tracker}, nil
}

func (r *Repository) Add(ctx context.Context, aggregate *courier.Courier) error {
	return r.saveWithTx(ctx, DomainToDTO(aggregate), true)
}

func (r *Repository) Update(ctx context.Context, aggregate *courier.Courier) error {
	return r.saveWithTx(ctx, DomainToDTO(aggregate), false)
}

func (r *Repository) saveWithTx(ctx context.Context, dto CourierDTO, isCreate bool) error {
	r.tracker.Track(DtoToDomain(dto)) // предположим, dto знает как получить доменную модель

	inTx := r.tracker.InTx()
	if !inTx {
		r.tracker.Begin(ctx)
	}
	tx := r.tracker.Tx()

	var err error
	session := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true})
	if isCreate {
		err = session.Create(&dto).Error
	} else {
		err = session.Save(&dto).Error
	}
	if err != nil {
		return err
	}

	if !inTx {
		if err = r.tracker.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, ID uuid.UUID) (*courier.Courier, error) {
	var dto CourierDTO
	result := r.txOrDb().WithContext(ctx).Preload(clause.Associations).First(&dto, ID)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, errs.NewObjectNotFoundError(ID.String(), nil)
		}
		return nil, result.Error
	}
	return DtoToDomain(dto), nil
}

func (r *Repository) GetAllFree(ctx context.Context) ([]*courier.Courier, error) {
	var dtos []CourierDTO

	err := r.txOrDb().WithContext(ctx).
		Preload(clause.Associations).
		Where(`
			NOT EXISTS (
				SELECT 1 FROM storage_places sp
				WHERE sp.courier_id = couriers.id AND sp.order_id IS NOT NULL
			)`).Find(&dtos).Error

	if err != nil {
		return nil, err
	}
	if len(dtos) == 0 {
		return nil, errs.NewObjectNotFoundError("Free couriers", nil)
	}

	aggregates := make([]*courier.Courier, len(dtos))
	for i := range dtos {
		aggregates[i] = DtoToDomain(dtos[i])
	}
	return aggregates, nil
}

func (r *Repository) txOrDb() *gorm.DB {
	if tx := r.tracker.Tx(); tx != nil {
		return tx
	}
	return r.tracker.Db()
}
