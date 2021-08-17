package staff

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
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

func TestNewPassportBearer(t *testing.T) {
	mustConfigViper()
	key := getSigningKey()

	pp, err := NewPassport(mockAccount, key)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Passport: %+v", pp)

	claims, err := ParsePassportClaims(pp.Token, key)
	if err != nil {
		t.Error(err)
	}

	t.Logf("Parsed PassportClaims: %+v", claims)
}
