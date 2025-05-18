package commands

import (
	"context"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
	"errors"
)

type MoveCouriersCommandHandler interface {
	Handle(context.Context, *MoveCouriersCommand) error
}

var _ MoveCouriersCommandHandler = &moveCouriersCommandHandler{}

type moveCouriersCommandHandler struct {
	unitOfWork ports.UnitOfWork
}

func NewMoveCouriersCommandHandler(
	unitOfWork ports.UnitOfWork) (MoveCouriersCommandHandler, error) {
	if unitOfWork == nil {
		return nil, errs.NewValueIsRequiredError("unitOfWork")
	}

	return &moveCouriersCommandHandler{
		unitOfWork: unitOfWork}, nil
}

func (ch *moveCouriersCommandHandler) Handle(ctx context.Context, command *MoveCouriersCommand) error {
	if command == nil {
		return errs.NewValueIsRequiredError("add address command")
	}

	// Восстановили
	assignedOrders, err := ch.unitOfWork.OrderRepository().GetAllInAssignedStatus(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return nil
		}
		return err
	}

	// Изменили и сохранили
	ch.unitOfWork.Begin(ctx)
	for _, assignedOrder := range assignedOrders {
		courier, err := ch.unitOfWork.CourierRepository().Get(ctx, *assignedOrder.CourierId())
		if err != nil {
			if errors.Is(err, errs.ErrObjectNotFound) {
				return nil
			}
			return err
		}

		err = courier.StepTowards(assignedOrder.Location())
		if err != nil {
			return err
		}

		if courier.Location().Equals(assignedOrder.Location()) {
			err := assignedOrder.Complete()
			if err != nil {
				return err
			}
			err = courier.CompleteOrder(assignedOrder)
			if err != nil {
				return err
			}
		}

		err = ch.unitOfWork.OrderRepository().Update(ctx, assignedOrder)
		if err != nil {
			return err
		}
		err = ch.unitOfWork.CourierRepository().Update(ctx, courier)
		if err != nil {
			return err
		}
	}
	err = ch.unitOfWork.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
