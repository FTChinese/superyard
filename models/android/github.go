package android

import "net/url"

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

// GitHubRelease is the published full release for the repository.
// https://developer.github.com/v3/repos/releases/#get-a-single-release
type GitHubRelease struct {
	ID          int64  `json:"id"`
	TagName     string `json:"tag_name"`
	Body        string `json:"body"`
	Draft       bool   `json:"draft"`
	CreateAt    string `json:"create_at"`
	PublishedAt string `json:"published_at"`
}
