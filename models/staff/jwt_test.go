package staff

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestJWTSign(t *testing.T) {
	mySigningKey := []byte("AllYourBase")

	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 86400*7,
		Id:        "stf_1234567890",
		IssuedAt:  time.Now().Unix(),
		Issuer:    "ftc_superyard",
		NotBefore: 0,
		Subject:   "",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		t.Error(err)
	}

	t.Log(ss)

	token, err = jwt.ParseWithClaims(ss, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySigningKey, nil
	})

	if err != nil {
		t.Error(err)
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		t.Logf("Is token valid: %t", claims.Valid() == nil)
	} else {
		t.Error(err)
	}
}
