package controller

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type HomeCtx struct {
	Year int
}

func Home(c echo.Context) error {
	c.Response().Header().Add("Cache-Control", "no-cache")
	c.Response().Header().Add("Cache-Control", "no-store")
	c.Response().Header().Add("Cache-Control", "must-revalidate")
	c.Response().Header().Add("Pragma", "no-cache")

	return c.Render(http.StatusOK, "home", HomeCtx{
		Year: time.Now().In(chrono.TZShanghai).Year(),
	})
}
