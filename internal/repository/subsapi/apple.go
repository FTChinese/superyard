package subsapi

import (
	"github.com/FTChinese/superyard/pkg/apple"
	"github.com/FTChinese/superyard/pkg/fetch"
	"net/http"
)

func (c Client) LinkIAP(link apple.LinkInput) (*http.Response, []error) {
	url := c.baseURL + "/apple/link"

	return fetch.New().
		Post(url).
		SetAuth(c.key).
		SendJSON(link).
		End()
}

// ListIAPSubs fetch a list of IAP subscriptions.
// The query string is forwarded as is.
// It does not have the `?` sign.
func (c Client) ListIAPSubs(query string) (*http.Response, []error) {
	url := c.baseURL + "/apple/subscription?" + query

	return fetch.New().
		Get(url).
		SetAuth(c.key).
		End()
}

func (c Client) LoadIAPSubs(id string) (*http.Response, []error) {
	url := c.baseURL + "/apple/subscription/" + id

	return fetch.New().
		Get(url).
		SetAuth(c.key).
		End()
}

func (c Client) RefreshIAPSubs(id string) (*http.Response, []error) {
	url := c.baseURL + "/apple/subscription/" + id

	return fetch.New().
		Patch(url).
		SetAuth(c.key).
		End()
}
