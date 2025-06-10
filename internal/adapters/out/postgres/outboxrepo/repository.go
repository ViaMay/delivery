package outboxrepo

import (
	"context"
	"delivery/internal/pkg/errs"
	"delivery/internal/pkg/outbox"
	"gorm.io/gorm"
)

type OutboxRepository interface {
	Update(ctx context.Context, message *outbox.Message) error
	GetNotPublishedMessages() ([]*outbox.Message, error)
}

var _ OutboxRepository = &repository{}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (OutboxRepository, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}

	return &repository{
		db: db,
	}, nil
}

func (r *repository) Update(ctx context.Context, message *outbox.Message) error {
	err := r.db.WithContext(ctx).Save(&message).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetNotPublishedMessages() ([]*outbox.Message, error) {
	var events []*outbox.Message
	result := r.db.
		Order("occurred_at_utc ASC").
		Limit(20).
		Where("processed_at_utc IS NULL").Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return events, nil
}
