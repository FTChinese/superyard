package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/staff"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	claimsCtxKey = "claims"
	keyUserID    = "X-User-Id"
	keyUnionID   = "X-Union-Id"
)

type AuthGuard struct {
	signingKey []byte
}

func NewAuthGuard(key []byte) AuthGuard {
	return AuthGuard{signingKey: key}
}

func (g AuthGuard) getPassportClaims(req *http.Request) (staff.PassportClaims, error) {
	authHeader := req.Header.Get("Authorization")
	ss, err := ParseBearer(authHeader)
	if err != nil {
		log.Printf("Error parsing Authorization header: %v", err)
		return staff.PassportClaims{}, err
	}

	claims, err := staff.ParsePassportClaims(ss, g.signingKey)
	if err != nil {
		log.Printf("Error parsing JWT %v", err)
		return staff.PassportClaims{}, err
	}

	return claims, nil
}

func (g AuthGuard) RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := g.getPassportClaims(c.Request())
		if err != nil {
			log.Printf("Error parsing JWT %v", err)
			return render.NewUnauthorized(err.Error())
		}

		c.Set(claimsCtxKey, claims)
		return next(c)
	}
}

func getPassportClaims(c echo.Context) staff.PassportClaims {
	return c.Get(claimsCtxKey).(staff.PassportClaims)
}

func RequireUserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := strings.TrimSpace(c.Request().Header.Get(keyUserID))

		if userID == "" {
			return render.NewUnauthorized("Missing X-User-Id")
		}

		return next(c)
	}
}

func getUserID(c echo.Context) string {
	return c.Request().Header.Get(keyUserID)
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
