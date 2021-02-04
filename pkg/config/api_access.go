package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

// AuthKeys is used to contain api access token set authorization header.
// Those keys are always comes in pair, one for development and one for production.
type AuthKeys struct {
	Dev  string `mapstructure:"dev"`
	Prod string `mapstructure:"prod"`
	name string
}

func (k AuthKeys) Validate() error {
	if k.Dev == "" || k.Prod == "" {
		return errors.New("dev or prod key found")
	}

	return nil
}

func (k AuthKeys) Pick(debug bool) string {
	log.Printf("Using %s for debug %t", k.name, debug)

	if debug {
		return k.Dev
	}

	return k.Prod
}

func LoadAuthKeys(name string) (AuthKeys, error) {
	var keys AuthKeys
	err := viper.UnmarshalKey(name, &keys)
	if err != nil {
		return keys, err
	}

	if err := keys.Validate(); err != nil {
		return keys, err
	}

	keys.name = name

	return keys, nil
}

func MustLoadAuthKeys(name string) AuthKeys {
	k, err := LoadAuthKeys(name)
	if err != nil {
		log.Fatalf("cannot get %s: %s", name, err.Error())
	}

	return k
}

// MustLoadAPIKey gets the API authorization key used by current app.
func MustLoadAPIKey() AuthKeys {
	return MustLoadAuthKeys("api_keys.superyard")
}
