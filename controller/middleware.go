package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/superyard/pkg/staff"
	"log"
	"net/http/httputil"
)

const claimsCtxKey = "claims"

// Guard holds various signing keys.
type Guard struct {
	JWT     string `mapstructure:"jwt_signing_key"`
	CSRF    string `mapstructure:"csrf_signing_key"`
	jwtKey  []byte
	csrfKey []byte
}

// NewGuard gets the keys from viper config file.
func NewGuard(name string) (Guard, error) {
	var guardKey Guard
	err := viper.UnmarshalKey(name, &guardKey)
	if err != nil {
		return guardKey, err
	}

	guardKey.jwtKey = []byte(guardKey.JWT)
	guardKey.csrfKey = []byte(guardKey.CSRF)

	return guardKey, nil
}

func MustNewGuard() Guard {
	k, err := NewGuard("web_app.superyard")
	if err != nil {
		log.Fatal(err)
	}

	return k
}

func (g Guard) createPassport(account staff.Account) (staff.PassportBearer, error) {
	return staff.NewPassportBearer(account, g.jwtKey)
}

func (g Guard) RequireLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		ss, err := ParseBearer(authHeader)
		if err != nil {
			log.Printf("Error parsing Authorization header: %v", err)
			return render.NewUnauthorized(err.Error())
		}

		claims, err := staff.ParsePassportClaims(ss, g.jwtKey)
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
