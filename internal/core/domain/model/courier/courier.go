package courier

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrNoFreeStoragePlace = errors.New("no free storage place")
)

type Courier struct {
	id               uuid.UUID
	name             string
	speed            int
	location         kernel.Location
	StoragePlaceList []*StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if speed <= 0 {
		return nil, errs.NewValueIsRequiredError("speed")
	}
	if err := location.IsValid(); err != nil {
		return nil, err
	}

	return &Courier{
		id:               uuid.New(),
		name:             name,
		speed:            speed,
		location:         location,
		StoragePlaceList: make([]*StoragePlace, 0),
	}, nil
}

func (c *Courier) AddStoragePlace(name string, totalVolume int) error {
	storagePlace, err := NewStoragePlace(name, totalVolume)
	if err != nil {
		return err
	}
	c.StoragePlaceList = append(c.StoragePlaceList, storagePlace)
	return nil

}

func (c *Courier) CanTakeOrder(order *order.Order) (bool, error) {
	if order == nil {
		return false, errs.NewValueIsRequiredError("order")
	}

	for _, storagePlace := range c.StoragePlaceList {
		canStore := storagePlace.CanStore(order.Volume())
		if canStore {
			return true, nil
		}
	}
	return false, nil
}

func (c *Courier) TakeOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}
	for _, storagePlace := range c.StoragePlaceList {
		if storagePlace.CanStore(order.Volume()) {
			return storagePlace.StoreOrder(order.Id(), order.Volume())
		}
	}
	return ErrNoFreeStoragePlace
}

func (c *Courier) CompleteOrder(order *order.Order) error {
	if order == nil {
		return errs.NewValueIsRequiredError("order")
	}
	for _, storagePlace := range c.StoragePlaceList {
		if storagePlace.OrderID() != nil && *storagePlace.OrderID() == order.Id() {
			storagePlace.RemoveOrder()
			return nil
		}
	}
	return nil
}

func (c *Courier) StepsTo(target kernel.Location) (float64, error) {
	if err := target.IsValid(); err != nil {
		return 0, err
	}
	distance, err := c.location.DistanceTo(target)
	if err != nil {
		return 0, err
	}

	time := float64(distance) / float64(c.speed)
	return time, err

}

func (c *Courier) StepTowards(target kernel.Location) error {
	if err := target.IsValid(); err != nil {
		return err
	}
	newX, newY := c.location.X(), c.location.Y()

	if newX < target.X() {
		newX++
	} else if newX > target.X() {
		newX--
	} else if newY < target.Y() {
		newY++
	} else if newY > target.Y() {
		newY--
	}

	newLoc, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLoc
	return nil
}

func (c *Courier) Equal(other *Courier) bool {
	if other == nil {
		return false
	}
	return c.id == other.id
}

func (c *Courier) ID() uuid.UUID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Speed() int {
	return c.speed
}

func (c *Courier) Location() kernel.Location {
	return c.location
}

func (c *Courier) StoragePlaces() []StoragePlace {
	res := make([]StoragePlace, len(c.StoragePlaceList))
	for i, storagePlace := range c.StoragePlaceList {
		res[i] = *storagePlace
	}
	return res
}
