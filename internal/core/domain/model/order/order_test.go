package order_test

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func invalidLocation() (kernel.Location, error) {
	return kernel.NewLocation(-999, 0)

}

var (
	testOrderID  uuid.UUID
	testVolume   int
	testLocation kernel.Location
)

func setup() {
	testOrderID = uuid.New()
	testVolume = 10
	testLocation, _ = kernel.NewLocation(5, 5)
}

func Test_NewOrder_Success(t *testing.T) {
	setup()
	o, err := order.NewOrder(testOrderID, testLocation, testVolume)

	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, testOrderID, o.Id())
	assert.Equal(t, testLocation, o.Location())
	assert.Equal(t, testVolume, o.Volume())
	assert.Equal(t, order.StatusCreated, o.Status())
	assert.Nil(t, o.CourierId())
}

func Test_NewOrder_EmptyId(t *testing.T) {
	setup()

	o, err := order.NewOrder(uuid.Nil, testLocation, testVolume)

	assert.Nil(t, o)
	assert.Error(t, err)
}

func Test_NewOrder_InvalidLocation(t *testing.T) {
	setup()
	loc, _ := invalidLocation()
	o, err := order.NewOrder(testOrderID, loc, testVolume)

	assert.Nil(t, o)
	assert.Error(t, err)
}

func Test_NewOrder_ZeroVolume(t *testing.T) {
	setup()

	o, err := order.NewOrder(testOrderID, testLocation, 0)

	assert.Nil(t, o)
	assert.Error(t, err)
}

func Test_Order_AssignCourier_Success(t *testing.T) {
	setup()

	o, _ := order.NewOrder(testOrderID, testLocation, testVolume)
	courierID := uuid.New()

	err := o.AssignCourier(courierID)

	assert.NoError(t, err)
	assert.NotNil(t, o.CourierId())
	assert.Equal(t, courierID, *o.CourierId())
	assert.Equal(t, order.StatusAssigned, o.Status())
}

func Test_Order_AssignCourier_WrongStatus(t *testing.T) {
	setup()

	o, _ := order.NewOrder(testOrderID, testLocation, testVolume)
	courierID := uuid.New()

	_ = o.AssignCourier(courierID)
	err := o.AssignCourier(uuid.New())

	assert.Error(t, err)
}

func Test_Order_Complete_Success(t *testing.T) {
	setup()

	o, _ := order.NewOrder(testOrderID, testLocation, testVolume)
	err := o.Complete()

	assert.Error(t, err)
}

func Test_Order_Complete_WrongStatus(t *testing.T) {
	setup()

	o, _ := order.NewOrder(testOrderID, testLocation, testVolume)
	err := o.Complete()

	assert.Error(t, err)
}

func Test_Order_Equals(t *testing.T) {
	setup()

	o1, _ := order.NewOrder(testOrderID, testLocation, testVolume)
	o2, _ := order.NewOrder(testOrderID, testLocation, testVolume)
	o3, _ := order.NewOrder(uuid.New(), testLocation, testVolume)

	assert.True(t, o1.Equals(o2))
	assert.False(t, o1.Equals(o3))
	assert.False(t, o1.Equals(nil))
}
