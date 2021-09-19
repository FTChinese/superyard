package b2bapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"net/http"
)

// ListOrders retrieves a list of orders.
// TODO: add filter for pagination, team_id, status.
func (c B2BClient) ListOrders() (*http.Response, error) {
	url := c.baseURL + pathOrders

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

func (c B2BClient) ConfirmOrder(id string) (*http.Response, error) {
	url := c.baseURL + pathOrderOf(id)

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
