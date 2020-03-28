package controller

import (
	"errors"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/employee"
	"strings"
)

// ParseAuthHeader extracts Authorization header.
// Authorization: Bearer 19c7d9016b68221cc60f00afca7c498c36c361e3
func ParseBearer(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("empty authorization header")
	}

	s := strings.SplitN(authHeader, " ", 2)

	bearerExists := (len(s) == 2) && (strings.ToLower(s[0]) == "bearer")

	if !bearerExists {
		return "", errors.New("bearer not found")
	}

	return s[1], nil
}

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

func CheckJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		ss, err := ParseBearer(authHeader)
		if err != nil {
			return render.NewUnauthorized(err.Error())
		}

		claims, err := employee.ParseJWT(ss)
		if err != nil {
			return render.NewUnauthorized(err.Error())
		}

		c.Set("claims", claims)
		return next(c)
	}
}

func getAccountClaims(c echo.Context) employee.AccountClaims {
	return c.Get("claims").(employee.AccountClaims)
}
