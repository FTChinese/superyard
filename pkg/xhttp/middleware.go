package xhttp

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"log"
	"net/http/httputil"
	"strings"
)

// NoCache set Cache-Control request header
func NoCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h := c.Response().Header()
		h.Add("Cache-Control", "no-cache")
		h.Add("Cache-Control", "no-store")
		h.Add("Cache-Control", "must-revalidate")
		h.Add("Pragma", "no-cache")
		return next(c)
	}
}

func DumpRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dump, err := httputil.DumpRequest(c.Request(), false)
		if err != nil {
			log.Print(err)
		}

		log.Printf(string(dump))

		return next(c)
	}
}

func RequireUserIDHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := strings.TrimSpace(c.Request().Header.Get(XUserID))

		if userID == "" {
			return render.NewUnauthorized("Missing X-User-Id")
		}

		return next(c)
	}
}

func RequireUserIDsQuery(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ftcID := strings.TrimSpace(c.QueryParam("ftc_id"))
		unionID := strings.TrimSpace(c.QueryParam("union_id"))

		if ftcID == "" && unionID == "" {
			return render.NewUnauthorized("Missing ftc_id or union_id in  query parameters")
		}

		return next(c)
	}
}
