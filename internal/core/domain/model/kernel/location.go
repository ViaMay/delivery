package kernel

import (
	"errors"
	"math"
	"math/rand"
)

var (
	ErrorInvalidCoordinate = errors.New("coordinate must be between 1 and 10 inclusive")
)

type Location struct {
	x int
	y int
}

func NewLocation(x, y int) (Location, error) {
	if x < 1 || x > 10 || y < 1 || y > 10 {
		return Location{}, ErrorInvalidCoordinate
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

func (l Location) DistanceTo(other Location) int {
	dx := int(math.Abs(float64(l.x - other.x)))
	dy := int(math.Abs(float64(l.y - other.y)))
	return dx + dy
}

func CreateRandom() (Location, error) {
	x := rand.Intn(10) + 1
	y := rand.Intn(10) + 1
	location, err := NewLocation(x, y)
	return location, err
}
