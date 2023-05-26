package subsapi

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
)

func (c Client) StripeActivePrices(refresh bool) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripePrices).
		AddQueryBool(queryKeyRefresh, refresh).
		String()

	log.Printf("List stripe active prices at %s", url)

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

// ListStripePrices retrieves all stripe prices with pagination.
func (c Client) ListStripePrices(query url.Values, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsStripe).
		AddPath("prices").
		String()

	log.Printf("List stripe prices at %s", url)

	resp, errs := fetch.
		New().
		Get(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		WithQuery(query).
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

func (c Client) UpdateStripePriceMeta(id string, body io.Reader, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripePrices).
		AddPath(id).
		String()

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) ActivateStripePrice(id string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsStripe).
		AddPath(id).
		AddPath("activate").
		String()

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) DeactivateStripePrice(id string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsStripe).
		AddPath(id).
		AddPath("deactivate").
		String()

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.HeaderStaffName(by)).
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
		SetHeader(xhttp.HeaderStaffName(by)).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
