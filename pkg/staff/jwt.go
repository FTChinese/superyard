package staff

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

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

// NewPassportClaims create a instance from staff's account.
func NewPassportClaims(a Account) PassportClaims {
	now := time.Now().Unix()

	return PassportClaims{
		StaffID:  a.ID.String,
		Username: a.UserName,
		Groups:   a.GroupMembers,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now + 86400*7,
			IssuedAt:  now,
			Issuer:    "com.ftchinese.superyard",
		},
	}
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

func (c PassportClaims) SignedString(key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return ss, nil
}

// PassportBearer contains a user's full account data
// plus the JSON Web Token and its expiration time.
type PassportBearer struct {
	Account
	ExpiresAt int64  `json:"expiresAt"`
	Token     string `json:"token"`
}

// NewPassportBearer creates a new PassportBearer for an account.
func NewPassportBearer(a Account, signingKey []byte) (PassportBearer, error) {

	claims := NewPassportClaims(a)

	ss, err := claims.SignedString(signingKey)

	if err != nil {
		return PassportBearer{}, err
	}

	return PassportBearer{
		Account:   a,
		ExpiresAt: claims.ExpiresAt,
		Token:     ss,
	}, nil
}
