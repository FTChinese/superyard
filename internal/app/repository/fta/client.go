package fta

import "github.com/FTChinese/superyard/pkg/config"

const (
	baseCmsApi = "/api/cms"
	baseTerms  = "/terms"
	pathTeams  = baseCmsApi + "/teams"
	pathOrders = baseCmsApi + "/orders"
)

type Client struct {
	key     string
	baseURL string
}

func NewClient(prod bool) Client {
	return Client{
		key:     config.MustFtaAPIKey().Pick(prod),
		baseURL: config.MustFtaBaseURL().Pick(prod),
	}
}
