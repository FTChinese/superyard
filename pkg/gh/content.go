package gh

import "encoding/base64"

// Content is the contents of a file or directory from GitHub
// See https://developer.github.com/v3/repos/contents/#get-contents
type Content struct {
	Encoding string `json:"encoding"` // This is always `base64`
	Name     string `json:"name"`     // File name
	Content  string `json:"content"`
}

// Decode decodes the raw content of a GitHub file.
func (c Content) Decode() (string, error) {
	data, err := base64.StdEncoding.DecodeString(c.Content)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
