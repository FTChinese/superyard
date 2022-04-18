package xhttp

import (
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/labstack/echo/v4"
)

func GetHeaderFtcID(c echo.Context) string {
	return c.Request().Header.Get(XUserID)
}

func GetQueryLive(c echo.Context) bool {
	return conv.DefaultTrue(c.QueryParam("live"))
}

func GetQueryRefresh(c echo.Context) bool {
	return conv.DefaultFalse(c.QueryParam("refresh"))
}

func BuildHeaderStaffName(n string) (string, string) {
	return XStaffName, n
}

func BuildHeaderFtcID(id string) (string, string) {
	return XUserID, id
}
