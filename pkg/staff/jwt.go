package staff

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// AccountClaims is a JWT custom claims containing only the
// essential fields of an account so that the signed string
// won't become too long while the backend can determine
// user's identity.
// After user logged in, the JWT is send to client as one
// of the JSON fields. The response body contains more fields
// than this claims so that client is able to show extra
// information on UI.
type AccountClaims struct {
	StaffID  string `json:"sid"`
	Username string `json:"name"`
	Groups   int64  `json:"grp"`
	jwt.StandardClaims
}

// SignedString create a JWT based on current claims.
func (c AccountClaims) SignedString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	ss, err := token.SignedString(signingKey)

	if err != nil {
		return ss, err
	}

	return ss, nil
}

// JWTAccount adds ExpiresAt so that client
// could check whether the login session is
// expired. It carries the Json Web Token.
type JWTAccount struct {
	Account
	ExpiresAt int64  `json:"expiresAt"`
	Token     string `json:"token"`
}

func NewJWTAccount(a Account) (JWTAccount, error) {
	expiresAt := time.Now().Unix() + 86400*7
	claims := AccountClaims{
		StaffID:  a.ID.String,
		Username: a.UserName,
		Groups:   a.GroupMembers,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
			Issuer:    "com.ftchinese.superyard",
		},
	}

	ss, err := claims.SignedString()
	if err != nil {
		return JWTAccount{}, err
	}

	return JWTAccount{
		Account:   a,
		ExpiresAt: expiresAt,
		Token:     ss,
	}, nil
}

func ParseJWT(ss string) (AccountClaims, error) {
	token, err := jwt.ParseWithClaims(
		ss,
		&AccountClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return signingKey, nil
		})

	if err != nil {
		log.Printf("Parsing JWT error: %v", err)
		return AccountClaims{}, err
	}

	log.Printf("Claims: %v", token.Claims)

	// NOTE: token.Claims is an interface, so it is a pointer, not a value type.
	if claims, ok := token.Claims.(*AccountClaims); ok {
		return *claims, nil
	}
	return AccountClaims{}, errors.New("wrong JWT claims")
}
