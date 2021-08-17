package config

import (
	"github.com/spf13/viper"
	"log"
)

// AppKey contains various signing or access keys for an app.
type AppKey struct {
	JWT     string `mapstructure:"jwt_signing_key"`
	CSRF    string `mapstructure:"csrf_signing_key"`
	jwtKey  []byte
	csrfKey []byte
}

func (a AppKey) GetJWTKey() []byte {
	return a.jwtKey
}

func (a AppKey) GetCSRFKey() []byte {
	return a.csrfKey
}

func GetAppKey(name string) (AppKey, error) {
	var appKey AppKey
	err := viper.UnmarshalKey(name, &appKey)
	if err != nil {
		return appKey, err
	}

	appKey.jwtKey = []byte(appKey.JWT)
	appKey.csrfKey = []byte(appKey.CSRF)

	return appKey, nil
}

func MustGetAppKey() AppKey {
	k, err := GetAppKey("web_app.superyard")
	if err != nil {
		log.Fatal(err)
	}

	return k
}
