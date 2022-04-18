package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"io"
	"net/http"
)

func (c Client) LoadStripeCoupon(id string, refresh bool) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathStripeCoupons).
		AddPath(id).
		AddQueryBool("refresh", refresh).
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

func (c Client) UpdateStripeCoupon(id string, body io.Reader, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsStripeCoupons).
		AddPath(id).
		String()

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

func (c Client) DeleteCoupon(id string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsStripeCoupons).
		AddPath(id).
		String()

	resp, errs := fetch.New().
		Delete(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
