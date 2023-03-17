package controller

import (
	"log"
	"net/http"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
)

type AuthGuard struct {
	signingKey []byte
}

func NewAuthGuard(key []byte) AuthGuard {
	return AuthGuard{signingKey: key}
}

func (g AuthGuard) getPassportClaims(req *http.Request) (user.PassportClaims, error) {
	authHeader := req.Header.Get("Authorization")
	ss, err := xhttp.ParseBearer(authHeader)
	if err != nil {
		log.Printf("Error parsing Authorization header: %v", err)
		return user.PassportClaims{}, err
	}

	claims, err := user.ParsePassportClaims(ss, g.signingKey)
	if err != nil {
		log.Printf("Error parsing JWT %v", err)
		return user.PassportClaims{}, err
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

		c.Set(xhttp.ClaimsCtxKey, claims)
		return next(c)
	}
}

func getPassportClaims(c echo.Context) user.PassportClaims {
	return c.Get(xhttp.ClaimsCtxKey).(user.PassportClaims)
}
