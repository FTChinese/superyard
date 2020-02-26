package apps

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"gitlab.com/ftchinese/superyard/models/android"
)

var request = gorequest.New()

const ghApiBase = "https://api.github.com"
const repoName = "ftc-android-kotlin"

var ghRepoUrl = fmt.Sprintf("%s/repos/FTChinese/%s", ghApiBase, repoName)

var ghLatestReleaseURL = fmt.Sprintf("%s/releases/latest", ghRepoUrl)

func ghReferenceURL(tag string) string {
	return fmt.Sprintf("%s/git/refs/tags/%s", ghRepoUrl, tag)
}

func gradleFileURL(tag string) string {
	return fmt.Sprintf("%s/contents/app/bubild.gradle?ref=%s", ghRepoUrl, tag)
}

func ghSingleReleaseURL(tag string) string {
	return fmt.Sprintf("%s/releases/tags/%s", ghRepoUrl, tag)
}

type GHRepo struct {
	client android.GitHubClient
}

func NewGHRepo(c android.GitHubClient) GHRepo {
	return GHRepo{
		client: c,
	}
}

func (g GHRepo) LatestRelease() (android.GitHubRelease, error) {
	_, body, err := request.Get(ghLatestReleaseURL).
		Set("User-Agent", "FTChinese").
		Query(g.client.Query()).
		EndBytes()

	if err != nil {
		return android.GitHubRelease{}, err[0]
	}

	var r android.GitHubRelease
	if err := json.Unmarshal(body, &r); err != nil {
		return android.GitHubRelease{}, err
	}

	return r, nil
}
