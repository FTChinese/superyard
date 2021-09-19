package b2bapi

import "github.com/FTChinese/superyard/pkg/config"

const (
	pathTeams  = "/teams"
	pathOrders = "/orders"
)

func pathTeamOf(id string) string {
	return pathTeams + "/" + id
}

func pathOrderOf(id string) string {
	return pathOrders + "/" + id
}

type B2BClient struct {
	key     string
	baseURL string
}

func NewClient(prod bool) B2BClient {
	return B2BClient{
		key:     config.MustFtaAPIKey().Pick(prod),
		baseURL: config.MustB2BBaseURL().Pick(prod),
	}
}
