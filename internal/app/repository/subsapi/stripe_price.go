package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"log"
	"net/http"
)

func (c Client) ListStripePrices(refresh bool) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		SetPath(pathStripePrices).
		AddQueryBool(queryKeyRefresh, refresh).
		String()

	log.Printf("List stripe prices at %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) LoadStripePrice(id string, refresh bool) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		SetPath(pathStripePrices).
		SetPath(id).
		SetQueryBool("refresh", refresh).
		String()

	log.Printf("Load a stripe prices at %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
