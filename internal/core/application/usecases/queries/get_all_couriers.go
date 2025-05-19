package queries

import (
	"delivery/internal/pkg/errs"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GetAllCouriersQuery struct{}

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string
	Location LocationResponse `gorm:"embedded;embeddedPrefix:location_"`
}

type GetAllCouriersQueryHandler interface {
	Handle(GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

type getAllCouriersQueryHandler struct {
	db *gorm.DB
}

func NewGetAllCouriersQueryHandler(db *gorm.DB) (GetAllCouriersQueryHandler, error) {
	if db == nil {
		return nil, errs.NewValueIsRequiredError("db")
	}
	return &getAllCouriersQueryHandler{db: db}, nil
}

func (q *getAllCouriersQueryHandler) Handle(_ GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	var couriers []CourierResponse
	result := q.db.Raw("SELECT id, name, location_x, location_y FROM couriers").Scan(&couriers)

	if result.Error != nil {
		return GetAllCouriersResponse{}, result.Error
	}

	return GetAllCouriersResponse{Couriers: couriers}, nil
}
