package b2bapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"net/http"
)

func (c B2BClient) LoadTeam(id string) (*http.Response, error) {
	url := c.baseURL + pathTeamOf(id)

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
