package controller

import (
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/labstack/echo/v4"
)

func getParamLive(c echo.Context) bool {
	return conv.DefaultTrue(c.QueryParam("live"))
}
