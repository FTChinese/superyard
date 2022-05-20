package ftaapi

import "github.com/FTChinese/superyard/pkg/config"

const (
	baseCmsApi = "/api/cms"
	baseTerms  = "/terms"
	pathTeams  = baseCmsApi + "/teams"
	pathOrders = baseCmsApi + "/orders"
)

type FtaClient struct {
	key     string
	baseURL string
}

func NewClient(prod bool) FtaClient {
	return FtaClient{
		key:     config.MustFtaAPIKey().Pick(prod),
		baseURL: config.MustFtaBaseURL().Pick(prod),
	}
}
