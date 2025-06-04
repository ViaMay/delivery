package ports

import (
	"context"
	"delivery/internal/pkg/ddd"
)

type OrderProducer interface {
	Publish(ctx context.Context, domainEvent ddd.DomainEvent) error
	Close() error
}
