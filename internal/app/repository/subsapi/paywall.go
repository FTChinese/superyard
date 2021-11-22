package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"io"
	"log"
	"net/http"
)

// RefreshPaywall cache bust paywall of both live and sandbox server.
// Deprecated
func (c Client) RefreshPaywall() (*http.Response, error) {
	go func() {
		c.refreshPaywall(true)
	}()

	return c.refreshPaywall(false)
}

func (c Client) RefreshFtcPaywall() (*http.Response, error) {
	url := c.baseURL + pathRefreshPaywall

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

func (c Client) RefreshStripePrices() (*http.Response, error) {
	url := c.baseURL + pathStripePrices

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

// refreshPaywall bust cache of paywall, either live version or
// sandbox version.
// Currently, both live/sandbox version exists on the save endpoint.
// In the future, live server only contains live daa, while sandbox server for sandbox data only.
// Deprecated.
func (c Client) refreshPaywall(sandbox bool) (*http.Response, error) {
	var url string
	if sandbox {
		url = c.sandboxBaseURL + pathRefreshPaywall
	} else {
		url = c.v3BaseUrl + pathRefreshPaywall
	}

	log.Printf("Refreshing paywall data at %s", url)

	resp, errs := fetch.New().Get(url).SetBearerAuth(c.key).End()
	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// LoadPaywall data from API. It always returns the live version.
func (c Client) LoadPaywall() (*http.Response, error) {
	url := c.baseURL + basePathPaywall

	resp, errs := fetch.New().Get(url).SetBearerAuth(c.key).End()
	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// CreatePrice creates a new price for a product.
// Input:
// createBy: string;
// tier: string;
// cycle: string;
// description?: string;
// liveMode: boolean;
// nickname?: string;
// price: number;
// productId: string;
func (c Client) CreatePrice(body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathProductPrices

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		Send(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// ListPriceOfProduct loads all prices under a product.
func (c Client) ListPriceOfProduct(productID string) (*http.Response, error) {
	url := c.baseURL + pathPricesOfProduct(productID)

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// ActivatePrice by id and returned and activated FtcPrice.
func (c Client) ActivatePrice(priceID string) (*http.Response, error) {
	url := c.baseURL + pathPriceOf(priceID)

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// RefreshPriceDiscounts update a price's discount list.
// Returns the updated FtcPrice.
func (c Client) RefreshPriceDiscounts(priceID string) (*http.Response, error) {
	url := c.baseURL + pathPriceOf(priceID)

	resp, errs := fetch.New().
		Patch(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// CreateDiscount for a price and returns the created one.
// Input:
// createdBy: string;
// description?: string;
// kind: introductory | promotion | retention | win_back
// percent: number;
// startUtc: string;
// endUtc: string;
// priceOff: number;
// priceId: string;
// recurring: boolean;
func (c Client) CreateDiscount(body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathPriceDiscounts

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		Send(body).End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// RemoveDiscount from a ftc price.
// Returns FtcPrice
func (c Client) RemoveDiscount(id string) (*http.Response, error) {
	url := c.baseURL + pathDiscountOf(id)

	resp, errs := fetch.New().
		Delete(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
