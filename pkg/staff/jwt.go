package staff

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

func NewStandardClaims(expiresAt int64) jwt.StandardClaims {
	return jwt.StandardClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  time.Now().Unix(),
		Issuer:    "com.ftchinese.superyard",
	}
}

// PassportClaims is a JWT custom claims to be signed as
// JSON Web Token. It contains only the
// essential fields of an account so that the signed string
// won't become too long while the backend can determine
// user's identity.
// After user logged in, the JWT is send to client as one
// of the JSON fields. The response body contains more fields
// than this claims so that client is able to show extra
// information on UI.
type PassportClaims struct {
	StaffID  string `json:"sid"`
	Username string `json:"name"`
	Groups   int64  `json:"grp"`
	jwt.StandardClaims
}

// Passport contains a user's full account data
// plus the JSON Web Token and its expiration time.
type Passport struct {
	Account
	ExpiresAt int64  `json:"expiresAt"`
	Token     string `json:"token"`
}

// NewPassport creates a new Passport for an account.
func NewPassport(a Account, signingKey []byte) (Passport, error) {

	claims := PassportClaims{
		StaffID:        a.ID.String,
		Username:       a.UserName,
		Groups:         a.GroupMembers,
		StandardClaims: NewStandardClaims(time.Now().Unix() * 86400 * 7),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)

	if err != nil {
		return Passport{}, err
	}

	return Passport{
		Account:   a,
		ExpiresAt: claims.ExpiresAt,
		Token:     ss,
	}, nil
}

// ParsePassportClaims parses a string to a PassportClaims
func ParsePassportClaims(ss string, key []byte) (PassportClaims, error) {
	token, err := jwt.ParseWithClaims(
		ss,
		&PassportClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return key, nil
		})

	if err != nil {
		log.Printf("Parsing JWT error: %v", err)
		return PassportClaims{}, err
	}

	// NOTE: token.Claims is an interface, so it is a pointer, not a value type.
	if claims, ok := token.Claims.(*PassportClaims); ok {
		return *claims, nil
	}
	return PassportClaims{}, errors.New("wrong JWT claims")
}
