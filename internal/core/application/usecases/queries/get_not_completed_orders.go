package queries

import (
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetNotCompletedOrdersQuery struct{}

type GetNotCompletedOrdersResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	ID       uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Location LocationResponse `gorm:"embedded;embeddedPrefix:location_"`
}

func (OrderResponse) TableName() string {
	return "orders"
}

type GetNotCompletedOrdersQueryHandler interface {
	Handle(GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error)
}

type getNotCompletedOrdersQueryHandler struct {
	db *gorm.DB
}

func NewGetNotCompletedOrdersQueryHandler(db *gorm.DB) (GetNotCompletedOrdersQueryHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &getNotCompletedOrdersQueryHandler{db: db}, nil
}

func (h *getNotCompletedOrdersQueryHandler) Handle(query GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error) {
	var orders []OrderResponse

	result := h.db.Raw(`
		SELECT id, location_x, location_y 
		FROM orders 
		WHERE status != ?`, order.StatusCompleted).Scan(&orders)

	if result.Error != nil {
		return GetNotCompletedOrdersResponse{}, result.Error
	}

	return GetNotCompletedOrdersResponse{Orders: orders}, nil
}
