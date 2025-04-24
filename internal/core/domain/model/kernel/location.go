package kernel

import (
	"delivery/internal/pkg/errs"
	"errors"
	"fmt"
	"math"
	"math/rand"
)

var (
	ErrorInvalidCoordinate = errors.New("coordinate must be between 1 and 10 inclusive")
)

const (
	minX int = 1
	maxX int = 10
	minY int = 1
	maxY int = 10
)

type Location struct {
	x int
	y int
}

func NewLocation(x, y int) (Location, error) {
	if x < minX || x > maxX {
		return Location{}, errs.NewValueIsOutOfRangeError("x", x, minX, maxX)
	}
	if y < minY || y > maxY {
		return Location{}, errs.NewValueIsOutOfRangeError("y", y, minY, maxY)
	}
	return Location{
		x: x,
		y: y,
	}, nil
}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) Equals(other Location) bool {
	return l.x == other.x && l.y == other.y
}

func (l Location) DistanceTo(other Location) (int, error) {
	if err := l.IsValid(); err != nil {
		return 0, fmt.Errorf("invalid origin location: %w", err)
	}
	if err := other.IsValid(); err != nil {
		return 0, fmt.Errorf("invalid target location: %w", err)
	}

	dx := int(math.Abs(float64(l.x - other.x)))
	dy := int(math.Abs(float64(l.y - other.y)))
	return dx + dy, nil
}

func (l Location) IsValid() error {
	if l.x < minX || l.x > maxX {
		return errs.NewValueIsOutOfRangeError("x", l.x, minX, maxX)
	}
	if l.y < minY || l.y > maxY {
		return errs.NewValueIsOutOfRangeError("y", l.y, minY, maxY)
	}
	return nil
}

func CreateRandom() (Location, error) {
	x := rand.Intn(maxX) + 1
	y := rand.Intn(maxY) + 1
	location, err := NewLocation(x, y)
	if err != nil {
		panic("CreateRandom(): failed to create a valid Location: " + err.Error())
	}
	return location, err
}
