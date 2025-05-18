package postgres

import (
	"context"
	"delivery/internal/pkg/uow"
	"gorm.io/gorm"
)

type Mapper[DTO any, Domain any] interface {
	ToDTO(Domain) DTO
	ToDomain(DTO) Domain
}

type BaseRepository[DTO any, Domain any] struct {
	tracker uow.Tracker
	mapper  Mapper[DTO, Domain]
}

func NewBaseRepository[DTO any, Domain any](tracker uow.Tracker, mapper Mapper[DTO, Domain]) *BaseRepository[DTO, Domain] {
	return &BaseRepository[DTO, Domain]{
		tracker: tracker,
		mapper:  mapper,
	}
}

func (r *BaseRepository[DTO, Domain]) Save(ctx context.Context, aggregate Domain, isCreate bool) error {
	r.tracker.Track(aggregate)
	dto := r.mapper.ToDTO(aggregate)

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
		return r.tracker.Commit(ctx)
	}
	return nil
}

func (r *BaseRepository[DTO, Domain]) GetByID(ctx context.Context, id any) (Domain, error) {
	var dto DTO
	err := r.txOrDb().WithContext(ctx).First(&dto, id).Error
	if err != nil {
		var zero Domain
		return zero, err
	}
	return r.mapper.ToDomain(dto), nil
}

func (r *BaseRepository[DTO, Domain]) txOrDb() *gorm.DB {
	if tx := r.tracker.Tx(); tx != nil {
		return tx
	}
	return r.tracker.Db()
}
