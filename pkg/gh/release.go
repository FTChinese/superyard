package gh

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// Release is the published full release for the repository.
// https://developer.github.com/v3/repos/releases/#get-a-single-release
type Release struct {
	ID          int64       `json:"id"`
	TagName     string      `json:"tag_name"`
	Body        null.String `json:"body"`
	Draft       bool        `json:"draft"`
	CreatedAt   chrono.Time `json:"created_at"`
	PublishedAt chrono.Time `json:"published_at"`
}
