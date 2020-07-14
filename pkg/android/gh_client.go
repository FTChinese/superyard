package android

import (
	"encoding/json"
	"errors"
	"github.com/FTChinese/go-rest/render"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

var request = gorequest.New()

const apkURL = "https://creatives.ftacademy.cn/minio/android/ftchinese-%s-ftc-release.apk"

const ghRepoURL = "https://api.github.com/repos/FTChinese/ftc-android-kotlin"

// GitHubClient contains the Github OAuth client id and secret
type GitHubClient struct {
	ID     string `mapstructure:"client_id"`
	Secret string `mapstructure:"client_secret"`
}

func NewGitHubClient(key string) (GitHubClient, error) {
	var c GitHubClient

	err := viper.UnmarshalKey(key, &c)
	if err != nil {
		return GitHubClient{}, err
	}

	return c, nil
}

func MustNewGitHubClient() GitHubClient {
	c, err := NewGitHubClient("oauth_client.gh_superyard")

	if err != nil {
		panic(err)
	}

	return c
}

// fetchRelease is used to fetch either the latest release
// or a single release of a specific tag.
func (c GitHubClient) fetchRelease(url string) (GitHubRelease, *render.ResponseError) {
	response, body, err := request.Get(url).
		Set("User-Agent", "FTChinese").
		SetBasicAuth(c.ID, c.Secret).
		EndBytes()

	if err != nil {
		return GitHubRelease{}, render.NewInternalError(err[0].Error())
	}

	// If response is not 200.
	if response != nil && response.StatusCode != 200 {

		return GitHubRelease{}, render.NewResponseError(response.StatusCode, response.Status)
	}

	var r GitHubRelease
	if err := json.Unmarshal(body, &r); err != nil {
		return GitHubRelease{}, render.NewBadRequest(err.Error())
	}

	return r, nil
}

// GetLatestRelease from https://api.github.com/repos/FTChinese/ftc-android-kotlin/releases/latest
func (c GitHubClient) GetLatestRelease() (GitHubRelease, *render.ResponseError) {
	return c.fetchRelease(ghRepoURL + "/releases/latest")
}

// GetSingleRelease from https://api.github.com/repos/FTChinese/ftc-android-kotlin/releases/tags/<tag>
func (c GitHubClient) GetSingleRelease(tag string) (GitHubRelease, *render.ResponseError) {
	return c.fetchRelease(ghRepoURL + "/releases/tags/" + tag)
}

func (c GitHubClient) GetRawContent(url string) (GitHubContent, *render.ResponseError) {
	r, body, err := request.
		Get(url).
		Set("User-Agent", "FTChinese").
		SetBasicAuth(c.ID, c.Secret).
		EndBytes()

	if err != nil {
		return GitHubContent{}, render.NewInternalError(err[0].Error())
	}

	if r != nil && r.StatusCode != 200 {
		return GitHubContent{}, render.NewResponseError(r.StatusCode, r.Status)
	}

	var content GitHubContent
	if err := json.Unmarshal(body, &c); err != nil {
		return GitHubContent{}, render.NewBadRequest(err.Error())
	}

	return content, nil
}

// GetGradleFile from https://api.github.com/repos/FTChinese/ftc-android-kotlin/contents/app/build.gradle?ref=<tag>
// The sole purpose of fetching the Gradle file is to extract
// versionCode.
func (c GitHubClient) GetGradleFile(tag string) (GitHubContent, *render.ResponseError) {
	return c.GetRawContent(ghRepoURL + "/contents/app/build.gradle?ref=" + tag)
}

// ParseVersionCode gets the value of versionCode field
// from a gradle file.
func ParseVersionCode(content string) (int64, error) {
	codeStr := extractVersionCode(content)

	if codeStr == "" {
		return 0, errors.New("versionCode missing in gradle file")
	}

	versionCode, err := strconv.ParseInt(codeStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return versionCode, nil
}

// extractVersionCode get the versionName field from gradle file.
func extractVersionCode(str string) string {
	lines := strings.Split(str, "\n")

	var target string
	for _, v := range lines {
		if strings.Contains(v, "versionCode") {
			target = v
			break
		}
	}

	if target == "" {
		return ""
	}
	parts := strings.Split(strings.TrimSpace(target), " ")
	return strings.TrimSpace(parts[len(parts)-1])
}
