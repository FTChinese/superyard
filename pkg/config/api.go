package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
)

// API is used to load configurations to access API.
// Those keys are always comes in pair, one for development and one for production.
type API struct {
	Dev  string `mapstructure:"dev"`
	Prod string `mapstructure:"prod"`
	name string
}

func (k API) Validate() error {
	if k.Dev == "" || k.Prod == "" {
		return errors.New("dev or prod key found")
	}

	return nil
}

func (k API) Pick(prod bool) string {
	log.Printf("Using %s for production %t", k.name, prod)

	if prod {
		return k.Prod
	}

	return k.Dev
}

func LoadAPIConfig(name string) (API, error) {
	var keys API
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

func MustLoadAPIConfig(name string) API {
	k, err := LoadAPIConfig(name)
	if err != nil {
		log.Fatalf("cannot get %s: %s", name, err.Error())
	}

	return k
}

// MustSubsAPIKey gets the API authorization key used by current app.
func MustSubsAPIKey() API {
	return MustLoadAPIConfig("api_keys.superyard")
}

func MustSubsAPISandboxBaseURL() API {
	return MustLoadAPIConfig("api_urls.sandbox")
}

func MustSubsAPIV3BaseURL() API {
	return MustLoadAPIConfig("api_urls.subs_v3")
}

// MustSubsAPIV4BaseURL loads the base url for api v4.
func MustSubsAPIV4BaseURL() API {
	return MustLoadAPIConfig("api_urls.subs_v4")
}

// MustSubsAPIv5BaseURL loads the base url for api v5
func MustSubsAPIv5BaseURL() API {
	return MustLoadAPIConfig("api_urls.subs_v5")
}

func MustSubsAPIv6BaseURL() API {
	return MustLoadAPIConfig("api_urls.subs_v6")
}

// MustFtaAPIKey loads key to access ftacademy's restful API.
func MustFtaAPIKey() API {
	return MustLoadAPIConfig("api_keys.ftacademy")
}

func MustFtaCmsURL() API {
	return MustLoadAPIConfig("api_urls.b2b_cms")
}
