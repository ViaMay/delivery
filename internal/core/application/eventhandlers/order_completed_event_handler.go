package eventhandlers

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"
)

type orderCompletedDomainEventHandler struct {
	orderProducer ports.OrderProducer
}

func NewOrderCompletedDomainEventHandler(
	orderProducer ports.OrderProducer) (ddd.EventHandler, error) {
	if orderProducer == nil {
		return nil, errs.NewValueIsRequiredError("orderProducer")
	}

	return &orderCompletedDomainEventHandler{orderProducer: orderProducer}, nil
}

func (eh *orderCompletedDomainEventHandler) Handle(ctx context.Context, domainEvent ddd.DomainEvent) error {
	err := eh.orderProducer.Publish(ctx, domainEvent)
	if err != nil {
		return err
	}
	return nil
}
