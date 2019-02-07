package subs

import (
	"strings"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/util"
)

// Banner is the content used on promotion banner
type Banner struct {
	CoverURL   string   `json:"coverUrl"`
	Heading    string   `json:"heading"`    // Required. Max 256 chars.
	SubHeading string   `json:"subHeading"` // Optional. Max 256 chars.
	Content    []string `json:"content"`    // Optional.
}

// Sanitize removes leading and trailing spaces.
func (b *Banner) Sanitize() {
	b.CoverURL = strings.TrimSpace(b.CoverURL)
	b.Heading = strings.TrimSpace(b.Heading)
	b.SubHeading = strings.TrimSpace(b.SubHeading)
}

// Validate validates input data for promotion banner.
func (b *Banner) Validate() *view.Reason {
	if r := util.OptionalMaxLen(b.CoverURL, 256, "coverUrl"); r != nil {
		return r
	}

	if r := util.RequireNotEmptyWithMax(b.Heading, 256, "heading"); r != nil {
		return r
	}

	return util.OptionalMaxLen(b.SubHeading, 256, "subHeading")
}

