package subsapi

import (
	"github.com/FTChinese/superyard/pkg/config"
)

type Client struct {
	key     string
	baseURL string
}

func NewClient(prod bool) Client {

	return Client{
		key:     config.MustLoadOAuthKey().Pick(prod),
		baseURL: config.MustSubsAPIv2BaseURL().Pick(prod),
	}
}
