package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DispatchService(t *testing.T) {
	// Arrange
	makeCourier := func(name string, x, y int) *courier.Courier {
		loc, _ := kernel.NewLocation(x, y)
		c, _ := courier.NewCourier(name, 1, loc)
		err := c.AddStoragePlace("Сумка", 10)
		if err != nil {
			return nil
		}
		return c
	}
	// Data
	courier1 := makeCourier("Pedestrian 1", 1, 1)
	courier2 := makeCourier("Pedestrian 2", 2, 2) // <- should win
	courier3 := makeCourier("Pedestrian 3", 3, 3)
	couriers := []*courier.Courier{courier1, courier2, courier3}

	orderLocation, err := kernel.NewLocation(2, 2)
	orderAggregate, err := order.NewOrder(uuid.New(), orderLocation, 5)

	dispatchService := NewDispatchService()

	// Act
	winner, err := dispatchService.Dispatch(orderAggregate, couriers)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, courier2, winner)
	assert.Equal(t, orderAggregate.Id(), *courier2.StoragePlaces()[0].OrderID())
	assert.Equal(t, order.StatusAssigned, orderAggregate.Status())
	assert.Equal(t, courier2.ID(), *orderAggregate.CourierId())
}
