package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"io"
	"net/http"
)

func (c Client) ListProduct() (*http.Response, error) {
	url := c.baseURL + pathProducts

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) CreateProduct(body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathProducts

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

func (c Client) LoadProduct(id string) (*http.Response, error) {
	url := c.baseURL + pathProductOf(id)

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) UpdateProduct(id string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + pathProductOf(id)

	resp, errs := fetch.New().
		Patch(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) ActivateProduct(id string) (*http.Response, error) {
	url := c.baseURL + pathActivateProductOf(id)

	resp, errs := fetch.New().
		Post(url).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
