package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"io"
	"net/http"
)

func (c Client) ListProduct(by string) (*http.Response, error) {
	url := c.baseURL + pathProducts

	resp, errs := fetch.New().
		Get(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) CreateProduct(body io.Reader, by string) (*http.Response, error) {
	url := c.baseURL + pathProducts

	resp, errs := fetch.New().
		Post(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) LoadProduct(id string, by string) (*http.Response, error) {
	url := c.baseURL + pathProductOf(id)

	resp, errs := fetch.New().
		Get(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) UpdateProduct(id string, body io.Reader, by string) (*http.Response, error) {
	url := c.baseURL + pathProductOf(id)

	resp, errs := fetch.New().
		Patch(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) ActivateProduct(id string, by string) (*http.Response, error) {
	url := c.baseURL + pathActivateProductOf(id)

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
