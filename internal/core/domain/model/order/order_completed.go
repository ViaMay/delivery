package order

import (
	"delivery/internal/pkg/ddd"
	"github.com/google/uuid"
	"reflect"
)

var _ ddd.DomainEvent = &CompletedDomainEvent{}

type CompletedDomainEvent struct {
	// base
	ID   uuid.UUID
	Name string

	// payload
	OrderID     uuid.UUID
	OrderStatus string

	isSet bool
}

func (e CompletedDomainEvent) GetID() uuid.UUID { return e.ID }

func (e CompletedDomainEvent) GetName() string {
	return e.Name
}

func NewCompletedDomainEvent(aggregate *Order) ddd.DomainEvent {
	completedDomainEvent := CompletedDomainEvent{
		ID: uuid.New(),

		OrderID:     aggregate.Id(),
		OrderStatus: aggregate.Status().String(),

		isSet: true,
	}
	completedDomainEvent.Name = reflect.TypeOf(completedDomainEvent).Name()
	return &completedDomainEvent
}

func NewEmptyCompletedDomainEvent() ddd.DomainEvent {
	completedDomainEvent := CompletedDomainEvent{}
	completedDomainEvent.Name = reflect.TypeOf(completedDomainEvent).Name()
	return &completedDomainEvent
}

func (e CompletedDomainEvent) IsEmpty() bool {
	return !e.isSet
}
