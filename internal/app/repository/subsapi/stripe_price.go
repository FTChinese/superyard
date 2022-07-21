package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"log"
	"net/http"
)

func (c Client) ListStripePrices(refresh bool) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripePrices).
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
		AddPath(pathStripePrices).
		AddPath(id).
		AddQueryBool("refresh", refresh).
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

// ListStripePriceCoupons for CMS, regardless of its current status.
func (c Client) ListStripePriceCoupons(priceID string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripePrices).
		AddPath(priceID).
		AddPath("coupons").
		String()

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
