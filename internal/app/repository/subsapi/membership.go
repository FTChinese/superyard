package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"io"
	"net/http"
	"net/url"
)

func (c Client) LoadMembership() (*http.Response, error) {
	to := fetch.
		NewURLBuilder(c.baseURL).
		AddPath(pathCmsMembership).
		String()

	resp, errs := fetch.
		New().
		Get(to).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) CreateMembership(body io.Reader, by string) (*http.Response, error) {
	to := fetch.
		NewURLBuilder(c.baseURL).
		AddPath(pathCmsMembership).
		String()

	resp, errs := fetch.
		New().
		Post(to).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) UpdateMembership(id string, body io.Reader, by string) (*http.Response, error) {
	to := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsMembership).
		AddPath(id).
		String()

	resp, errs := fetch.
		New().
		Patch(to).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) DeleteMembership(id string, body io.Reader, by string) (*http.Response, error) {
	to := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsMembership).
		AddPath(id).
		String()

	resp, errs := fetch.
		New().
		Delete(to).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		StreamJSON(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) ListSnapshot(query url.Values, by string) (*http.Response, error) {
	to := fetch.NewURLBuilder(c.baseURL).AddPath(pathCmsSnapshots).String()

	resp, errs := fetch.
		New().
		Get(to).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		WithQuery(query).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
