package config

import "github.com/spf13/viper"

type ApiBaseURL struct {
	ReaderV1    string `mapstructure:"reader_v1"`
	SubsV1      string `mapstructure:"subscription_v1"`
	SubsSandbox string `mapstructure:"sub_sandbox"`
	ContentV1   string `mapstructure:"content_v1"`
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
