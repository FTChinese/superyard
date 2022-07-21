package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"io"
	"log"
	"net/http"
)

// RefreshFtcPaywall hit the refresh endpoint.
// Deprecated. Used only for backward compatible.
func (c Client) RefreshFtcPaywall() (*http.Response, error) {
	url := c.baseURL + pathRefreshPaywall

	log.Printf("Refresh paywall at %s", url)

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

// LoadPaywall data from API. It always returns the live version.
func (c Client) LoadPaywall(refresh bool) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(rootPathPaywall).
		AddQueryBool(queryKeyRefresh, refresh).
		String()

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()
	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) CreatePaywallBanner(body io.Reader, by string) (*http.Response, error) {
	url := c.baseURL + pathPaywallBanner

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		SetHeader(xhttp.HeaderStaffName(by)).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) CreatePaywallPromoBanner(body io.Reader, by string) (*http.Response, error) {
	url := c.baseURL + pathPaywallPromo

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		SetHeader(xhttp.HeaderStaffName(by)).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) DropPaywallPromo(by string) (*http.Response, error) {
	url := c.baseURL + pathPaywallPromo

	resp, errs := fetch.New().
		Delete(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
