package subsapi

import (
	apple2 "github.com/FTChinese/superyard/internal/pkg/apple"
	"github.com/FTChinese/superyard/pkg/fetch"
	"net/http"
)

func (c Client) LinkIAP(link apple2.LinkInput) (*http.Response, []error) {
	url := c.sandboxBaseURL + "/apple/link"

	return fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		SendJSON(link).
		End()
}

func (c Client) UnlinkIAP(link apple2.LinkInput) (*http.Response, []error) {
	url := c.sandboxBaseURL + "/apple/unlink"

	return fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		SendJSON(link).End()
}

// ListIAPSubs fetch a list of IAP subscriptions.
// The query string is forwarded as is.
// It does not have the `?` sign.
func (c Client) ListIAPSubs(userID string, query string) (*http.Response, []error) {
	url := c.sandboxBaseURL + "/apple/subs?" + query

	return fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader("X-User-Id", userID).
		End()
}

func (c Client) LoadIAPSubs(id string) (*http.Response, []error) {
	url := c.sandboxBaseURL + "/apple/subs/" + id

	return fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()
}

func (c Client) RefreshIAPSubs(id string) (*http.Response, []error) {
	url := c.sandboxBaseURL + "/apple/subs/" + id

	return fetch.New().
		Patch(url).
		SetBearerAuth(c.key).
		End()
}
