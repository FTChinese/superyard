package android

import (
	"encoding/base64"
	"fmt"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// GitHubContent is the contents of a file or directory from GitHub
// See https://developer.github.com/v3/repos/contents/#get-contents
type GitHubContent struct {
	Encoding string `json:"encoding"` // This is always `base64`
	Name     string `json:"name"`     // File name
	Content  string `json:"content"`
}

// Decode decodes the raw content of a GitHub file.
func (c GitHubContent) Decode() (string, error) {
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

// FtcRelease turns GitHubRelease into FTC's Release.
func (r GitHubRelease) FtcRelease(versionCode int64) Release {
	return Release{
		ReleaseInput: ReleaseInput{
			VersionName: r.TagName,
			VersionCode: versionCode,
			Body:        r.Body,
			ApkURL:      fmt.Sprintf(apkURL, r.TagName),
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.PublishedAt,
	}
}
