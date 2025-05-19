package commands

import (
	"context"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateCourierCommandHandler interface {
	Handle(context.Context, *CreateCourierCommand) error
}

var _ CreateCourierCommandHandler = &createCourierCommandHandler{}

type createCourierCommandHandler struct {
	unitOfWork ports.UnitOfWork
}

func NewCreateCourierCommandHandler(
	unitOfWork ports.UnitOfWork,
) (CreateCourierCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &createCourierCommandHandler{
		unitOfWork: unitOfWork,
	}, nil
}

func (ch *createCourierCommandHandler) Handle(ctx context.Context, command *CreateCourierCommand) error {
	if command == nil {
		return errs.NewValueIsRequiredError("create courier command")
	}

	location, _ := kernel.CreateRandom()

	courierAggregate, err := courier.NewCourier(command.Name, command.Speed, location)
	if err != nil {
		return err
	}

	return ch.unitOfWork.CourierRepository().Add(ctx, courierAggregate)
}
