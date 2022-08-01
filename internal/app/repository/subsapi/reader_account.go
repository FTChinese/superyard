package subsapi

import (
	"github.com/FTChinese/superyard/internal/pkg/sandbox"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"net/http"
)

func (c Client) SignUp(p sandbox.SignUpParams, header http.Header) (fetch.Response, error) {
	url := c.baseURL + pathEmailSignUp

	resp, errs := fetch.
		New().
		Post(url).
		WithHeader(header).
		SetBearerAuth(c.key).
		SendJSON(p).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}

func (c Client) LoadFtcAccount(id string) (*http.Response, error) {
	url := c.baseURL + rootPathAccount

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader(xhttp.XUserID, id).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) LoadWxAccount(unionID string) (*http.Response, error) {
	url := c.baseURL + pathWxAccount

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader(xhttp.XUnionID, unionID).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) DeleteFtcAccount(a sandbox.TestAccount) (*http.Response, error) {
	url := c.baseURL + rootPathAccount

	resp, errs := fetch.New().
		Delete(url).
		SetBearerAuth(c.key).
		SetHeader(xhttp.XUserID, a.FtcID).
		SendJSON(sandbox.SignUpParams{
			Email:    a.Email,
			Password: a.ClearPassword,
		}).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) LoadFtcAddress(id string) (*http.Response, error) {
	url := c.baseURL + pathAddress

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader(xhttp.XUserID, id).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) LoadFtcProfile(id string) (*http.Response, error) {
	url := c.baseURL + pathProfile

	resp, errs := fetch.New().
		Get(url).
		SetBearerAuth(c.key).
		SetHeader(xhttp.XUserID, id).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
