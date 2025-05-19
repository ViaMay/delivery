package commands

import (
	"delivery/internal/pkg/errs"
)

type CreateCourierCommand struct {
	Name  string
	Speed int
}

func NewCreateCourierCommand(name string, speed int) (*CreateCourierCommand, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if speed <= 0 {
		return nil, errs.NewValueIsRequiredError("speed")
	}

	return &CreateCourierCommand{
		Name:  name,
		Speed: speed,
	}, nil
}
