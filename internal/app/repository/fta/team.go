package fta

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"net/http"
)

func (c Client) LoadTeam(id string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).AddPath(pathTeams).AddPath(id).String()

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
