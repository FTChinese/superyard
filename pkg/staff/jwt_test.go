package staff

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mustConfigViper() {
	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}
}

func getSigningKey() []byte {
	keyStr := viper.GetString("web_app.superyard.jwt_signing_key")

	return []byte(keyStr)
}

var mockClaims = PassportClaims{
	StaffID:        "stf_X3UccHoHqHMW",
	Username:       "RonaldCrawford",
	Groups:         2,
	StandardClaims: jwt.StandardClaims{},
}

func TestNewPassportClaims(t *testing.T) {
	claims := NewPassportClaims(mockAccount)

	assert.Equal(t, claims, mockClaims)
}

func TestPassportClaims_SignedString(t *testing.T) {
	mustConfigViper()
	key := getSigningKey()

	assert.NotEmpty(t, key)

	claims := NewPassportClaims(mockAccount)

	ss, err := claims.SignedString(key)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Signed string: %s", ss)
}

func TestNewPassportBearer(t *testing.T) {
	mustConfigViper()
	key := getSigningKey()

	pp, err := NewPassportBearer(mockAccount, key)
	if err != nil {
		t.Error(err)
	}

	t.Logf("PassportBearer: %+v", pp)

	claims, err := ParsePassportClaims(pp.Token, key)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Parsed PassportClaims: %+v", claims)
}
