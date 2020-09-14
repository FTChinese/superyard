package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"log"
	"net/http"
)

func (c Client) RefreshPaywall() (*http.Response, error) {
	url := c.baseURL + "/paywall/__refresh"

	log.Printf("Refreshing paywall data at %s", url)

	resp, errs := fetch.New().Get(url).SetAuth(c.key).End()
	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
