package controller

import (
	"github.com/guregu/null"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ftchinese/superyard/pkg/staff"
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

var mockAccount = staff.Account{
	ID:           null.StringFrom("stf_X3UccHoHqHMW"),
	UserName:     "RonaldCrawford",
	Email:        "kMeyer@Talane.info",
	DisplayName:  null.StringFrom("DebraAdams"),
	Department:   null.StringFrom("tech"),
	GroupMembers: 2,
	IsActive:     true,
}

func TestNewGuard(t *testing.T) {
	mustConfigViper()

	g := MustNewGuard()

	assert.NotEmpty(t, g.JWT)
	assert.NotEmpty(t, g.jwtKey)

	assert.Equal(t, g.jwtKey, getSigningKey())
}

func TestGuard_createPassport(t *testing.T) {
	mustConfigViper()

	g := MustNewGuard()

	pb, err := g.createPassport(mockAccount)

	if err != nil {
		t.Error(err)
	}

	t.Logf("Passport: %+v", pb)

	ss, err := staff.NewPassportClaims(mockAccount).
		SignedString(g.jwtKey)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, pb.Token, ss)

	t.Logf("Signed string: %s", ss)
}
