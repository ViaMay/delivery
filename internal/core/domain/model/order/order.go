package order

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

type Order struct {
	*ddd.BaseAggregate[uuid.UUID]
	courierId *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	if err := location.IsValid(); err != nil {
		return nil, err
	}
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}
	return &Order{
		BaseAggregate: ddd.NewBaseAggregate[uuid.UUID](orderID),
		location:      location,
		volume:        volume,
		status:        StatusCreated,
	}, nil
}

func (o *Order) AssignCourier(courierId uuid.UUID) error {
	if o.status != StatusCreated {
		return errors.New("order must be in Created status to assign courier")
	}
	o.courierId = &courierId
	o.status = StatusAssigned
	return nil
}

func (o *Order) Complete() error {
	if o.status != StatusAssigned {
		return errors.New("only assigned orders can be completed")
	}
	o.status = StatusCompleted

	// Публикуем доменное событие
	o.BaseAggregate.RaiseDomainEvent(NewCompletedDomainEvent(o))

	return nil
}

func (o *Order) Equals(other *Order) bool {
	if other == nil {
		return false
	}
	return o.ID() == other.ID()
}

func (o *Order) Id() uuid.UUID {
	return o.ID()
}

func (o *Order) CourierId() *uuid.UUID {
	return o.courierId
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Volume() int {
	return o.volume
}

func (o *Order) Status() Status {
	return o.status
}

func RestoreOrder(id uuid.UUID, courierID *uuid.UUID, location kernel.Location, volume int, status Status) *Order {
	return &Order{
		BaseAggregate: ddd.NewBaseAggregate(id),
		courierId:     courierID,
		location:      location,
		volume:        volume,
		status:        status,
	}
}
