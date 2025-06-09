package ports

import (
	"context"
	"delivery/internal/pkg/uow"
)

type UnitOfWork interface {
	uow.Tracker
	Begin(ctx context.Context)
	Commit(ctx context.Context) error
	CourierRepository() CourierRepository
	OrderRepository() OrderRepository
}
