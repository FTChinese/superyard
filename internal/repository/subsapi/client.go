package subsapi

import (
	"github.com/FTChinese/superyard/pkg/config"
)

type Client struct {
	key     string
	baseURL string
}

func NewClient(debug bool) Client {

	return Client{
		key:     config.MustLoadAPIKey().Pick(debug),
		baseURL: config.MustApiBaseURLs().GetSubsV1(debug),
	}
}
