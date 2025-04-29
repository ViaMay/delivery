package courier

import (
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/google/uuid"
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderId     *uuid.UUID
}

var (
	ErrCannotStoreOrder = errors.New("cannot store order: place is not empty or volume too large")
)

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if totalVolume <= 0 {
		return nil, errs.NewValueIsRequiredError("totalVolume")
	}
	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
		orderId:     nil,
	}, nil
}

func (s *StoragePlace) ID() uuid.UUID {
	return s.id
}

func (s *StoragePlace) Name() string {
	return s.name
}

func (s *StoragePlace) TotalVolume() int {
	return s.totalVolume
}

func (s *StoragePlace) OrderID() *uuid.UUID {
	return s.orderId
}

func (s *StoragePlace) IsEmpty() bool {
	return s.orderId == nil
}

func (s *StoragePlace) CanStore(orderVolume int) bool {
	if orderVolume <= 0 {
		return false
	}
	return s.IsEmpty() && orderVolume <= s.totalVolume
}

func (s *StoragePlace) StoreOrder(orderId uuid.UUID, orderVolume int) error {
	if orderId == uuid.Nil {
		return errs.NewValueIsRequiredError("orderId")
	}
	if orderVolume <= 0 {
		return errs.NewValueIsRequiredError("volume")
	}
	if !s.CanStore(orderVolume) {
		return ErrCannotStoreOrder
	}
	s.orderId = &orderId
	return nil
}

func (s *StoragePlace) RemoveOrder() {
	s.orderId = nil
}

func (s *StoragePlace) Equals(other *StoragePlace) bool {
	if other == nil {
		return false
	}
	return s.id == other.id
}
