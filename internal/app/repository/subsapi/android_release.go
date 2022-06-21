package subsapi

import (
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"io"
	"net/http"
)

// ListAndroidRelease loads all prices under a product.
func (c Client) ListAndroidRelease(rawQuery string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathAndroidReleases).
		SetRawQuery(rawQuery).
		String()

	resp, errs := fetch.
		New().
		Get(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) AndroidReleaseOf(versionName string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathAndroidReleases).
		AddPath(versionName).
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

func (c Client) CreateAndroidRelease(body io.Reader, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsAndroid).
		String()

	resp, errs := fetch.
		New().
		Post(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		Stream(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) UpdateAndroidRelease(versionName string, body io.Reader, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsAndroid).
		AddPath(versionName).
		String()

	resp, errs := fetch.
		New().
		Patch(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		Stream(body).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}

func (c Client) DeleteAndroidRelease(versionName string, by string) (*http.Response, error) {
	url := fetch.NewURLBuilder(c.baseURL).
		AddPath(pathCmsAndroid).
		AddPath(versionName).
		String()

	resp, errs := fetch.
		New().
		Delete(url).
		SetHeader(xhttp.BuildHeaderStaffName(by)).
		SetBearerAuth(c.key).
		End()

	if errs != nil {
		return nil, errs[0]
	}

	return resp, nil
}
