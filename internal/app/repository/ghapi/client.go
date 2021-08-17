package ghapi

import "github.com/spf13/viper"

const androidBaseURL = "https://api.github.com/repos/FTChinese/ftc-android-kotlin"

var userAgent = map[string]string{
	"User-Agent": "FTChinese",
}

type Client struct {
	ID     string `mapstructure:"client_id"`
	Secret string `mapstructure:"client_secret"`
}

func NewClient() (Client, error) {
	var c Client

	err := viper.UnmarshalKey("oauth_client.gh_superyard", &c)
	if err != nil {
		return Client{}, err
	}

	return c, nil
}

func MustNewClient() Client {
	c, err := NewClient()
	if err != nil {
		panic(err)
	}

	return c
}
