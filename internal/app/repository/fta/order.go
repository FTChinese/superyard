package fta

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"io"
	"net/http"
)

// ListOrders retrieves a list of orders.
func (c Client) ListOrders(rawQuery string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathOrders).
		SetRawQuery(rawQuery).
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

func (c Client) LoadOrder(id string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathOrders).
		AddPath(id).
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

func (c Client) ConfirmOrder(id string, body io.Reader) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathOrders).
		AddPath(id).
		String()

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
