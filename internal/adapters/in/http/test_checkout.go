package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) TestCheckout(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
