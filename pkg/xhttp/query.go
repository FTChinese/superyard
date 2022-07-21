package xhttp

import (
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetQueryLive(c echo.Context) bool {
	return conv.DefaultTrue(c.QueryParam("live"))
}

func GetQueryRefresh(c echo.Context) bool {
	return conv.DefaultFalse(c.QueryParam("refresh"))
}

func GetFtcID(c echo.Context) string {
	return c.QueryParam("ftc_id")
}

func GetHeaderWxID(c echo.Context) string {
	return c.QueryParam("union_id")
}

var decoder = schema.NewDecoder()

func DecodeForm(v interface{}, req *http.Request) error {
	decoder.IgnoreUnknownKeys(true)

	if err := req.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(v, req.Form); err != nil {
		return err
	}

	return nil
}
