package courier_test

import (
	"delivery/internal/core/domain/model/courier"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewStoragePlace_Success(t *testing.T) {
	name := "bag"
	totalVolume := 10

	sp, err := courier.NewStoragePlace(name, totalVolume)

	assert.NoError(t, err)
	assert.NotNil(t, sp)
	assert.Equal(t, name, sp.Name())
	assert.Equal(t, totalVolume, sp.TotalVolume())
	assert.True(t, sp.IsEmpty())
	assert.Nil(t, sp.OrderID())
}

func Test_StoragePlace_NewStoragePlace_Empty_Name(t *testing.T) {
	sp, err := courier.NewStoragePlace("", 10)
	assert.Nil(t, sp)
	assert.Contains(t, err.Error(), "name")
}

func Test_StoragePlace_NewStoragePlace_Invalid_Volume(t *testing.T) {
	sp, err := courier.NewStoragePlace("bag", 0)
	assert.Nil(t, sp)
	assert.Contains(t, err.Error(), "totalVolume")
}

func Test_StoragePlace_IsEmpty(t *testing.T) {
	sp, _ := courier.NewStoragePlace("bag", 10)
	assert.True(t, sp.IsEmpty())

	orderID := uuid.New()
	_ = sp.StoreOrder(orderID, 5)

	assert.False(t, sp.IsEmpty())
}

func Test_StoragePlace_CanStore(t *testing.T) {
	sp, _ := courier.NewStoragePlace("bag", 10)

	assert.True(t, sp.CanStore(5))
	assert.False(t, sp.CanStore(15))

	orderID := uuid.New()
	_ = sp.StoreOrder(orderID, 5)

	assert.False(t, sp.CanStore(5))
}

func Test_StoragePlace_StoreOrder_Success(t *testing.T) {
	sp, _ := courier.NewStoragePlace("Рюкзак", 10)
	orderID := uuid.New()

	err := sp.StoreOrder(orderID, 5)

	assert.NoError(t, err)
	assert.Equal(t, orderID, *sp.OrderID())
	assert.False(t, sp.IsEmpty())
}

func Test_StoragePlace_StoreOrder_Fail_NilOrderID(t *testing.T) {
	sp, _ := courier.NewStoragePlace("bag", 10)

	err := sp.StoreOrder(uuid.Nil, 5)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "orderId")
}

func Test_StoragePlace_StoreOrder_Fail_ZeroVolume(t *testing.T) {
	sp, _ := courier.NewStoragePlace("bag", 10)

	err := sp.StoreOrder(uuid.New(), 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "volume")
}

func Test_StoragePlace_StoreOrder_Fail_CannotStore(t *testing.T) {
	sp, _ := courier.NewStoragePlace("bag", 5)

	// Volume too large
	err := sp.StoreOrder(uuid.New(), 10)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, courier.ErrCannotStoreOrder))

	// Occupied
	orderID := uuid.New()
	_ = sp.StoreOrder(orderID, 5)

	err = sp.StoreOrder(uuid.New(), 3)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, courier.ErrCannotStoreOrder))
}

func Test_StoragePlace_RemoveOrder(t *testing.T) {
	sp, _ := courier.NewStoragePlace("bag", 10)
	orderID := uuid.New()
	_ = sp.StoreOrder(orderID, 5)

	assert.False(t, sp.IsEmpty())

	sp.RemoveOrder()
	assert.True(t, sp.IsEmpty())
	assert.Nil(t, sp.OrderID())
}

func Test_StoragePlace_Equals(t *testing.T) {
	sp1, _ := courier.NewStoragePlace("bag", 10)
	assert.True(t, sp1.Equals(sp1))

	sp2, _ := courier.NewStoragePlace("bag2", 15)

	assert.False(t, sp1.Equals(sp2))
}
