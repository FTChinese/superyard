package promo

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"strings"
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
func (b *Banner) Validate() *render.ValidationError {
	ie := validator.New("coverUrl").MaxLen(256).Validate(b.CoverURL)
	if ie != nil {
		return ie
	}

	ie = validator.New("heading").Required().MaxLen(256).Validate(b.Heading)
	if ie != nil {
		return ie
	}

	return validator.New("subHeading").MaxLen(256).Validate(b.SubHeading)
}
