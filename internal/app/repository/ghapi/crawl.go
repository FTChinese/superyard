package ghapi

import (
	"encoding/json"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/gh"
	"net/http"
	"net/url"
)

func (c Client) BuildFetch(url string) *fetch.Fetch {
	return fetch.New().
		Get(url).
		SetHeaderMap(userAgent).
		SetBasicAuth(c.ID, c.Secret)
}

func (c Client) Crawl(url string) (*http.Response, []byte, []error) {
	return c.BuildFetch(url).
		EndBytes()
}

func (c Client) getRelease(url string) (gh.Release, *render.ResponseError) {
	resp, body, errs := c.Crawl(url)

	if errs != nil {
		return gh.Release{}, render.NewInternalError(errs[0].Error())
	}

	if resp != nil && resp.StatusCode != 200 {
		return gh.Release{}, render.NewResponseError(resp.StatusCode, resp.Status)
	}

	var r gh.Release
	if err := json.Unmarshal(body, &r); err != nil {
		return gh.Release{}, render.NewBadRequest(err.Error())
	}

	return r, nil
}

// GetLatestRelease from https://api.github.com/repos/FTChinese/ftc-android-kotlin/releases/latest
func (c Client) GetLatestRelease(baseURL string) (gh.Release, *render.ResponseError) {
	return c.getRelease(baseURL + "/releases/latest")
}

// GetSingleRelease from https://api.github.com/repos/FTChinese/ftc-android-kotlin/releases/tags/<tag>
func (c Client) GetSingleRelease(baseURL, tag string) (gh.Release, *render.ResponseError) {
	return c.getRelease(baseURL + "/releases/tags/" + tag)
}

func (c Client) GetRawContent(url string, query url.Values) (gh.Content, *render.ResponseError) {
	f := c.BuildFetch(url)
	if query != nil {
		f.SetQuery(query)
	}

	resp, body, errs := f.EndBytes()

	if errs != nil {
		return gh.Content{}, render.NewInternalError(errs[0].Error())
	}

	if resp != nil && resp.StatusCode != 200 {
		return gh.Content{}, render.NewResponseError(resp.StatusCode, resp.Status)
	}

	var content gh.Content
	if err := json.Unmarshal(body, &content); err != nil {
		return gh.Content{}, render.NewBadRequest(err.Error())
	}

	return content, nil
}