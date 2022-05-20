package fta

import "github.com/FTChinese/superyard/pkg/config"

const (
	pathTeams  = "/teams"
	pathOrders = "/orders"
)

type Client struct {
	key     string
	baseURL string
}

func NewClient(prod bool) Client {
	return Client{
		key:     config.MustFtaAPIKey().Pick(prod),
		baseURL: config.MustFtaCmsURL().Pick(prod),
	}
}
