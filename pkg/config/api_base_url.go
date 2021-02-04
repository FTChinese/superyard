package config

import "github.com/spf13/viper"

type ApiBaseURL struct {
	SubsV1 string `mapstructure:"subscription_v1"`
}

func ApiBaseURLs() (ApiBaseURL, error) {
	var a ApiBaseURL
	err := viper.UnmarshalKey("api_url", &a)
	if err != nil {
		return a, err
	}

	return a, nil
}

func MustApiBaseURLs() ApiBaseURL {
	a, err := ApiBaseURLs()
	if err != nil {
		panic(err)
	}

	return a
}

func (u ApiBaseURL) GetSubsV1(debug bool) string {
	if debug {
		return "http://localhost:8200"

	}

	return u.SubsV1
}
