package services

import (
	"delivery/internal/core/domain/model/courier"
	orderpkg "delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"errors"
	"math"
)

type DispatchService interface {
	Dispatch(order *orderpkg.Order, couriers []*courier.Courier) (*courier.Courier, error)
}

type dispatchService struct{}

func NewDispatchService() DispatchService {
	return &dispatchService{}
}

var (
	ErrSuitableCourierWasNotFound = errors.New("suitable courier was not found")
	ErrOrderAlreadyAssigned       = errors.New("order is already assigned")
)

func (s *dispatchService) Dispatch(order *orderpkg.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if order == nil {
		return nil, errs.NewValueIsRequiredError("order")
	}
	if order.Status() != orderpkg.StatusCreated {
		return nil, ErrOrderAlreadyAssigned
	}
	if couriers == nil || len(couriers) == 0 {
		return nil, errs.NewValueIsRequiredError("couriers")
	}

	bestCourier, err := s.findBestCourier(order, couriers)
	if err != nil {
		return nil, err
	}

	if err := bestCourier.TakeOrder(order); err != nil {
		return nil, err
	}
	if err := order.AssignCourier(bestCourier.ID()); err != nil {
		return nil, err
	}
	return bestCourier, nil
}

func (s *dispatchService) findBestCourier(order *orderpkg.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	availableCouriers, err := s.filterAvailableCouriers(order, couriers)
	if len(availableCouriers) == 0 {
		return nil, ErrSuitableCourierWasNotFound
	}

	best, err := s.findFastestCourier(order, availableCouriers)
	if err != nil {
		return nil, err
	}
	return best, nil
}

func (s *dispatchService) filterAvailableCouriers(order *orderpkg.Order, couriers []*courier.Courier) ([]*courier.Courier, error) {
	var result []*courier.Courier
	for _, c := range couriers {
		canTake, _ := c.CanTakeOrder(order)
		if canTake {
			result = append(result, c)
		}
	}
	return result, nil
}

func (s *dispatchService) findFastestCourier(order *orderpkg.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	var (
		best    *courier.Courier
		minTime = math.MaxFloat64
	)
	location := order.Location()
	for _, c := range couriers {
		time, err := c.StepsTo(location)
		if err != nil {
			return nil, err
		}
		if time < minTime {
			minTime = time
			best = c
		}
	}
	return best, nil
}
