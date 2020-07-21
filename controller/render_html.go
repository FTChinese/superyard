package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Home(c echo.Context) error {
	c.Response().Header().Add("Cache-Control", "no-cache")
	c.Response().Header().Add("Cache-Control", "no-store")
	c.Response().Header().Add("Cache-Control", "must-revalidate")
	c.Response().Header().Add("Pragma", "no-cache")

	return c.Render(http.StatusOK, "home", nil)
}
