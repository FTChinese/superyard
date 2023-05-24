package subsapi

import (
	"io"
	"net/http"

	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
)

// ListPriceOfProduct loads all prices under a product.
func (c Client) ListPriceOfProduct(productID string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathPrices).
		AddQuery(queryKeyProductID, productID).
		String()

	resp, errs := fetch.New().
		Get(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// CreatePrice creates a new price for a product.
func (c Client) CreatePrice(body io.Reader, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).AddPath(pathPrices).String()

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		Stream(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) LoadFtcPrice(id string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).AddPath(pathPrices).AddPath(id).String()

	resp, errs := fetch.New().
		Get(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

// ActivatePrice by id and returned and activated FtcPrice.
func (c Client) ActivatePrice(priceID string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathPrices).
		AddPath(priceID).
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

func (c Client) UpdatePrice(id string, body io.Reader, by string) (*http.Response, error) {
	url := pathPriceOf(c.baseURL, id)

	resp, errs := fetch.New().
		Patch(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		Stream(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) ArchivePrice(id string, by string) (*http.Response, error) {
	url := pathPriceOf(c.baseURL, id)

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

// RefreshPriceDiscounts update a price's discount list.
// Returns the updated FtcPrice.
func (c Client) RefreshPriceDiscounts(priceID string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathPrices).
		AddPath(priceID).
		AddPath("discounts").
		String()

	resp, errs := fetch.New().
		Patch(url).
		SetHeader(xhttp.HeaderStaffName(by)).
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
	url := fetch.
		NewURLBuilder(c.baseURL).
		AddPath(pathPriceDiscounts).
		String()

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.HeaderStaffName(by)).
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
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathPriceDiscounts).
		AddPath(id).
		String()

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
