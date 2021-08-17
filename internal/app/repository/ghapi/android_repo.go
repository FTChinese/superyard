package ghapi

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/gh"
	"net/url"
)

// GetAndroidLatestRelease from https://api.github.com/repos/FTChinese/ftc-android-kotlin/releases/latest
func (c Client) GetAndroidLatestRelease() (gh.Release, *render.ResponseError) {
	return c.GetLatestRelease(androidBaseURL)
}

// GetAndroidRelease from https://api.github.com/repos/FTChinese/ftc-android-kotlin/releases/tags/<tag>
func (c Client) GetAndroidRelease(tag string) (gh.Release, *render.ResponseError) {
	return c.GetSingleRelease(androidBaseURL, tag)
}

// GetGradleFile from https://api.github.com/repos/FTChinese/ftc-android-kotlin/contents/app/build.gradle?ref=<tag>
// The sole purpose of fetching the Gradle file is to extract
// versionCode.
func (c Client) GetAndroidGradleFile(tag string) (gh.Content, *render.ResponseError) {
	return c.GetRawContent(
		androidBaseURL+"/contents/app/build.gradle",
		url.Values{
			"ref": []string{tag},
		},
	)
}
