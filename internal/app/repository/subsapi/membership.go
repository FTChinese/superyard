package subsapi

import (
	"errors"
	"github.com/FTChinese/superyard/pkg/fetch"
	"io"
	"net/http"
)

func (c Client) LoadMembership() (*http.Response, error) {
	url := c.baseURL + rootPathMember

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

func (c Client) CreateMembership(body io.Reader) (*http.Response, error) {
	url := c.baseURL + rootPathMember

	resp, errs := fetch.
		New().
		Post(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) UpdateMembership(body io.Reader) (*http.Response, error) {
	url := c.baseURL + rootPathMember

	resp, errs := fetch.
		New().
		Patch(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) DeleteMembership(body io.Reader) (*http.Response, error) {
	url := c.baseURL + rootPathMember

	resp, errs := fetch.
		New().
		Delete(url).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) ListSnapshot() (*http.Response, error) {
	return nil, errors.New("not implemented")
}
