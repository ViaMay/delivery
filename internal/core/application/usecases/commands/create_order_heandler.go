package commands

import (
	"context"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateOrderCommandHandler interface {
	Handle(context.Context, *CreateOrderCommand) error
}

var _ CreateOrderCommandHandler = &createOrderCommandHandler{}

type createOrderCommandHandler struct {
	unitOfWork ports.UnitOfWork
	geoClient  ports.GeoClient
}

func NewCreateOrderCommandHandler(unitOfWork ports.UnitOfWork, geoClient ports.GeoClient) (CreateOrderCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	if geoClient == nil {
		return nil, errs.NewValueIsRequiredError("geoClient")
	}

	return &createOrderCommandHandler{
		unitOfWork: unitOfWork,
		geoClient:  geoClient}, nil
}

func (ch *createOrderCommandHandler) Handle(ctx context.Context, command *CreateOrderCommand) error {
	if command == nil {
		return errs.NewValueIsRequiredError("create order command")
	}

	existingOrder, err := ch.unitOfWork.OrderRepository().Get(ctx, command.OrderID)
	if err != nil {
		return err
	}
	if existingOrder != nil {
		return nil
	}

	location, err := ch.geoClient.GetGeolocation(ctx, command.Street)
	if err != nil {
		return err
	}

	newOrder, err := order.NewOrder(command.OrderID, location, command.Volume)
	if err != nil {
		return err
	}

	return ch.unitOfWork.OrderRepository().Add(ctx, newOrder)
}
