package xhttp

import (
	"errors"
	"github.com/labstack/echo/v4"
	"strings"
)

func GetFtcID(c echo.Context) string {
	return c.Request().Header.Get(XUserID)
}

func HeaderStaffName(n string) (string, string) {
	return XStaffName, n
}

func HeaderFtcID(id string) (string, string) {
	return XUserID, id
}

// ParseBearer extracts Authorization header.
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
