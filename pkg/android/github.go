package android

import (
	"encoding/base64"
	"fmt"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"net/url"
)

const apkURL = "https://creatives.ftacademy.cn/minio/android/ftchinese-%s-ftc-release.apk"

// GitHubClient contains the Github OAuth client id and secret
type GitHubClient struct {
	ID     string `mapstructure:"client_id"`
	Secret string `mapstructure:"client_secret"`
}

func (c GitHubClient) Query() string {
	v := url.Values{}
	v.Set("client_id", c.ID)
	v.Set("client_secret", c.Secret)

	return v.Encode()
}

// GitHubContent is the contents of a file or directory from GitHub
// See https://developer.github.com/v3/repos/contents/#get-contents
type GitHubContent struct {
	Encoding string `json:"encoding"` // This is always `base64`
	Name     string `json:"name"`     // File name
	Content  string `json:"content"`
}

func (c GitHubContent) GetContent() (string, error) {
	data, err := base64.StdEncoding.DecodeString(c.Content)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// GitHubRelease is the published full release for the repository.
// https://developer.github.com/v3/repos/releases/#get-a-single-release
type GitHubRelease struct {
	ID          int64       `json:"id"`
	TagName     string      `json:"tag_name"`
	Body        null.String `json:"body"`
	Draft       bool        `json:"draft"`
	CreatedAt   chrono.Time `json:"created_at"`
	PublishedAt chrono.Time `json:"published_at"`
}

func (r GitHubRelease) FtcRelease(versionCode int64) Release {
	return Release{
		VersionName: r.TagName,
		VersionCode: versionCode,
		Body:        r.Body,
		ApkURL:      fmt.Sprintf(apkURL, r.TagName),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.PublishedAt,
	}
}
