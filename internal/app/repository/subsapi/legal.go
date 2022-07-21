package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"io"
	"net/http"
)

func (c Client) LoadLegalDoc(id string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathLegal).
		AddPath(id).
		String()

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

func (c Client) ListLegalDocs(rawQuery string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsLegal).
		SetRawQuery(rawQuery).
		String()

	resp, errs := fetch.
		New().
		Get(url).
		SetHeader(xhttp.HeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) CreateLegalDoc(body io.Reader, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsLegal).
		String()

	resp, errs := fetch.
		New().
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

func (c Client) UpdateLegalDoc(id string, body io.Reader, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsLegal).
		AddPath(id).
		String()

	resp, errs := fetch.
		New().
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

func (c Client) PublishLegalDoc(id string, body io.Reader, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsLegal).
		AddPath(id).
		AddPath("publish").
		String()

	resp, errs := fetch.
		New().
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
