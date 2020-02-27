package apps

import (
	"encoding/json"
	"fmt"
	"github.com/FTChinese/go-rest/render"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"gitlab.com/ftchinese/superyard/models/android"
	"strconv"
	"strings"
)

var request = gorequest.New()

const ghApiBase = "https://api.github.com"
const repoName = "ftc-android-kotlin"

var ghRepoUrl = fmt.Sprintf("%s/repos/FTChinese/%s", ghApiBase, repoName)

var ghLatestReleaseURL = fmt.Sprintf("%s/releases/latest", ghRepoUrl)

// The url to check if a tag exists.
func ghReferenceURL(tag string) string {
	return fmt.Sprintf("%s/git/refs/tags/%s", ghRepoUrl, tag)
}

func gradleFileURL(tag string) string {
	return fmt.Sprintf("%s/contents/app/build.gradle?ref=%s", ghRepoUrl, tag)
}

func ghSingleReleaseURL(tag string) string {
	return fmt.Sprintf("%s/releases/tags/%s", ghRepoUrl, tag)
}

func MustGetGH() android.GitHubClient {
	var c android.GitHubClient
	if err := viper.UnmarshalKey("oauth_client.gh_superyard", &c); err != nil {
		panic(err)
	}

	return c
}

type GHRepo struct {
	client android.GitHubClient
}

func NewGHRepo() GHRepo {
	return GHRepo{
		client: MustGetGH(),
	}
}

func (g GHRepo) fetchRelease(url string) (android.GitHubRelease, *render.ResponseError) {
	response, body, err := request.Get(url).
		Set("User-Agent", "FTChinese").
		Query(g.client.Query()).
		EndBytes()

	if err != nil {
		logger.WithField("trace", "GHRepo.fetchRelease").Error(err)
		return android.GitHubRelease{}, render.NewInternalError(err[0].Error())
	}

	if response != nil && response.StatusCode != 200 {
		logger.WithField("trace", "GHRepo.fetchRelease").Printf("Response status code: %d", response.StatusCode)

		return android.GitHubRelease{}, render.NewResponseError(response.StatusCode, response.Status)
	}

	var r android.GitHubRelease
	if err := json.Unmarshal(body, &r); err != nil {
		return android.GitHubRelease{}, render.NewBadRequest(err.Error())
	}

	return r, nil
}

func (g GHRepo) LatestRelease() (android.GitHubRelease, *render.ResponseError) {
	return g.fetchRelease(ghLatestReleaseURL)
}

func (g GHRepo) SingleRelease(tag string) (android.GitHubRelease, *render.ResponseError) {
	return g.fetchRelease(ghSingleReleaseURL(tag))
}

func (g GHRepo) GradleFile(tag string) (android.GitHubContent, *render.ResponseError) {
	r, body, err := request.Get(gradleFileURL(tag)).
		Set("User-Agent", "FTChinese").
		Query(g.client.Query()).
		EndBytes()

	if err != nil {
		return android.GitHubContent{}, render.NewInternalError(err[0].Error())
	}

	if r != nil && r.StatusCode != 200 {
		return android.GitHubContent{}, render.NewResponseError(r.StatusCode, r.Status)
	}

	var c android.GitHubContent
	if err := json.Unmarshal(body, &c); err != nil {
		return android.GitHubContent{}, render.NewBadRequest(err.Error())
	}

	return c, nil
}

func (g GHRepo) GetVersionCode(tag string) (int64, *render.ResponseError) {
	c, rErr := g.GradleFile(tag)

	if rErr != nil {
		return 0, rErr
	}

	fileContent, err := c.GetContent()
	if err != nil {
		return 0, render.NewBadRequest(err.Error())
	}

	codeStr := ExtractVersionCode(fileContent)

	if codeStr == "" {
		return 0, nil
	}

	versionCode, err := strconv.ParseInt(codeStr, 10, 64)
	if err != nil {
		return 0, nil
	}

	return versionCode, nil
}

// ExtractVersionCode get the versionName field from gradle file.
func ExtractVersionCode(str string) string {
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
