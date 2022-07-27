package subsapi

import (
	"github.com/FTChinese/superyard/internal/pkg/sandbox"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"net/http"
)

func (c Client) SignUp(p sandbox.SignUpParams, header http.Header) (fetch.Response, error) {

	resp, errs := fetch.
		New().
		Post(pathEmailSignUp).
		WithHeader(header).
		SetBearerAuth(c.key).
		SendJSON(p).
		EndBlob()

	if errs != nil {
		return fetch.Response{}, errs[0]
	}

	return resp, nil
}

func (c Client) LoadReader(id string) (*http.Response, error) {
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

func (c Client) DeleteReader(a sandbox.TestAccount) (*http.Response, error) {
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
