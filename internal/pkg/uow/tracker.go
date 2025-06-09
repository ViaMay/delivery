package uow

import (
	"context"
	"delivery/internal/pkg/ddd"
	"gorm.io/gorm"
)

type Tracker interface {
	Tx() *gorm.DB
	Db() *gorm.DB
	InTx() bool
	Track(agg ddd.AggregateRoot)
	Begin(ctx context.Context)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
