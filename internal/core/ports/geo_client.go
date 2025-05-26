package ports

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
)

type GeoClient interface {
	GetGeolocation(ctx context.Context, street string) (kernel.Location, error)
}
