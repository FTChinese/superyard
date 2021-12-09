package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"io"
	"net/http"
)

// ListPriceOfProduct loads all prices under a product.
func (c Client) ListPriceOfProduct(productID string, by string) (*http.Response, error) {
	url := c.baseURL + pathPrices

	resp, errs := fetch.New().
		Get(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		SetQuery(queryKeyProductID, productID).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// CreatePrice creates a new price for a product.
func (c Client) CreatePrice(body io.Reader, by string) (*http.Response, error) {
	url := c.baseURL + pathPrices

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		Stream(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// ActivatePrice by id and returned and activated FtcPrice.
func (c Client) ActivatePrice(priceID string, by string) (*http.Response, error) {
	url := c.baseURL + pathActivatePriceOf(priceID)

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) UpdatePrice(id string, body io.Reader, by string) (*http.Response, error) {
	url := c.baseURL + pathPriceOf(id)

	resp, errs := fetch.New().
		Patch(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		Stream(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// RefreshPriceDiscounts update a price's discount list.
// Returns the updated FtcPrice.
func (c Client) RefreshPriceDiscounts(priceID string, by string) (*http.Response, error) {
	url := c.baseURL + pathRefreshOffersOfPrice(priceID)

	resp, errs := fetch.New().
		Patch(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) ArchivePrice(id string, by string) (*http.Response, error) {
	url := c.baseURL + pathPriceOf(id)

	resp, errs := fetch.New().
		Delete(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
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
func (c Client) CreateDiscount(body io.Reader, by string) (*http.Response, error) {
	url := c.baseURL + pathPriceDiscounts

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		Stream(body).End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// RemoveDiscount from a ftc price.
// Returns FtcPrice
func (c Client) RemoveDiscount(id string, by string) (*http.Response, error) {
	url := c.baseURL + pathDiscountOf(id)

	resp, errs := fetch.New().
		Delete(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
