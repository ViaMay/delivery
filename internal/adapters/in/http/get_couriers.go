package http

import (
	"delivery/internal/adapters/in/http/problems"
	"delivery/internal/core/application/usecases/queries"
	"delivery/internal/generated/servers"
	"delivery/internal/pkg/errs"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetCouriers(c echo.Context) error {
	query, err := queries.NewGetAllCouriersQuery()
	if err != nil {
		return problems.NewBadRequest(err.Error())
	}

	response, err := s.getAllCouriersQueryHandler.Handle(query)
	if err != nil {
		if errors.Is(err, errs.ErrObjectNotFound) {
			return c.JSON(http.StatusNotFound, problems.NewNotFound(err.Error()))
		}
	}

	var couriers []servers.Courier
	for _, courier := range response.Couriers {
		location := servers.Location{
			X: courier.Location.X,
			Y: courier.Location.Y,
		}

		var courier = servers.Courier{
			Id:       courier.ID,
			Name:     courier.Name,
			Location: location,
		}
		couriers = append(couriers, courier)
	}
	return c.JSON(http.StatusOK, couriers)
}
