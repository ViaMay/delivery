package commands

import (
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
)

type CreateOrderCommand struct {
	OrderID uuid.UUID
	Street  string
	Volume  int
}

func NewCreateOrderCommand(orderID uuid.UUID, street string, volume int) (*CreateOrderCommand, error) {
	if street == "" {
		return nil, errs.NewValueIsRequiredError("street")
	}
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}

	return &CreateOrderCommand{
		OrderID: uuid.New(),
		Street:  street,
		Volume:  volume,
	}, nil
}
