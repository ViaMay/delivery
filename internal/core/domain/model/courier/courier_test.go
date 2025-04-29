package courier_test

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func validLocation() kernel.Location {
	loc, _ := kernel.NewLocation(1, 1)
	return loc
}

func targetLocation() kernel.Location {
	loc, _ := kernel.NewLocation(4, 4)
	return loc
}

func Test_NewCourier_Success(t *testing.T) {
	c, err := courier.NewCourier("CourierName", 3, validLocation())
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "CourierName", c.Name())
	assert.Equal(t, 3, c.Speed())
}

func Test_NewCourier_EmptyName(t *testing.T) {
	c, err := courier.NewCourier("", 3, validLocation())
	assert.Error(t, err)
	assert.Nil(t, c)
}

func Test_NewCourier_ZeroSpeed(t *testing.T) {
	c, err := courier.NewCourier("Courier", 0, validLocation())
	assert.Error(t, err)
	assert.Nil(t, c)
}

func Test_AddStoragePlace_Success(t *testing.T) {
	c, _ := courier.NewCourier("CourierName", 3, validLocation())
	err := c.AddStoragePlace("MainStorage", 10)
	assert.NoError(t, err)
	assert.Len(t, c.StoragePlaces(), 1)
}

func Test_CanTakeOrder_Success(t *testing.T) {
	c, _ := courier.NewCourier("CourierName", 3, validLocation())
	_ = c.AddStoragePlace("MainStorage", 10)

	o, _ := order.NewOrder(uuid.New(), targetLocation(), 5)

	canTake, err := c.CanTakeOrder(o)
	assert.NoError(t, err)
	assert.True(t, canTake)
}

func Test_CanTakeOrder_NoFreeSpace(t *testing.T) {
	c, _ := courier.NewCourier("CourierName", 3, validLocation())
	_ = c.AddStoragePlace("SmallStorage", 3) // маленький объём

	o, _ := order.NewOrder(uuid.New(), targetLocation(), 5) // заказ больше объёма

	canTake, err := c.CanTakeOrder(o)
	assert.NoError(t, err)
	assert.False(t, canTake)
}

func Test_TakeOrder_Success(t *testing.T) {
	c, _ := courier.NewCourier("CourierName", 3, validLocation())
	_ = c.AddStoragePlace("Storage", 10)

	o, _ := order.NewOrder(uuid.New(), targetLocation(), 5)

	err := c.TakeOrder(o)
	assert.NoError(t, err)
}

func Test_TakeOrder_NoStoragePlace(t *testing.T) {
	c, _ := courier.NewCourier("CourierName", 3, validLocation())
	_ = c.AddStoragePlace("SmallStorage", 3)

	o, _ := order.NewOrder(uuid.New(), targetLocation(), 5)

	err := c.TakeOrder(o)
	assert.ErrorIs(t, err, courier.ErrNoFreeStoragePlace)
}

func Test_CompleteOrder_Success(t *testing.T) {
	c, _ := courier.NewCourier("CourierName", 3, validLocation())
	_ = c.AddStoragePlace("Storage", 10)

	o, _ := order.NewOrder(uuid.New(), targetLocation(), 5)
	_ = c.TakeOrder(o)

	err := c.CompleteOrder(o)
	assert.NoError(t, err)
}

func Test_StepsTo(t *testing.T) {
	c, _ := courier.NewCourier("CourierName", 3, validLocation())

	steps, err := c.StepsTo(targetLocation())
	assert.NoError(t, err)
	assert.Equal(t, 2.0, steps)
}

func Test_StepTowards(t *testing.T) {
	c, _ := courier.NewCourier("CourierName", 3, validLocation())

	err := c.StepTowards(targetLocation())
	assert.NoError(t, err)

	assert.NotEqual(t, validLocation(), c.Location())
}

func Test_Courier_Equal(t *testing.T) {
	loc := validLocation()
	c1, _ := courier.NewCourier("Courier1", 3, loc)
	c2, _ := courier.NewCourier("Courier2", 3, loc)
	assert.True(t, c1.Equal(c1))
	assert.False(t, c1.Equal(c2))
	assert.False(t, c1.Equal(nil))
}
