package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"log"
	"net/http"
)

func (c Client) RefreshPaywall() (*http.Response, error) {
	url := c.baseURL + "/paywall/__refresh"

	log.Printf("Refreshing paywall data at %s", url)

	return fetch.NewRequest().Get(url).SetAuth(c.key).End()
}
