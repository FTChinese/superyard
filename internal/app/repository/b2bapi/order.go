package b2bapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"io"
	"net/http"
)

// ListOrders retrieves a list of orders.
func (c B2BClient) ListOrders(rawQuery string) (*http.Response, error) {
	url := c.baseURL + pathOrders + "?" + rawQuery

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c B2BClient) LoadOrder(id string) (*http.Response, error) {
	url := c.baseURL + pathOrderOf(id)

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c B2BClient) ConfirmOrder(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathOrderOf(id)

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
