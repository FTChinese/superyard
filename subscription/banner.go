package subscription

import (
	"encoding/json"
	"strings"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Banner is the content used on promotion banner
type Banner struct {
	Heading    string   `json:"heading"`    // Required. Max 256 chars.
	SubHeading string   `json:"subHeading"` // Optional. Max 256 chars.
	Content    []string `json:"content"`    // Optional.
}

// Sanitize removes leading and trailing spaces.
func (b *Banner) Sanitize() {
	b.Heading = strings.TrimSpace(b.Heading)
	b.SubHeading = strings.TrimSpace(b.SubHeading)
}

// Validate validates input data for promotion banner.
func (b *Banner) Validate() *util.Reason {
	if r := util.RequireNotEmptyWithMax(b.Heading, 256, "heading"); r != nil {
		return r
	}

	return util.OptionalMaxLen(b.SubHeading, 256, "subHeading")
}

// SaveBanner sets the banner content for a promotion.
// It is also used to edit banner content.
func (env Env) SaveBanner(id int64, banner Banner) error {
	query := `
	UPDATE premium.promotion_schedule
	SET banner = ?
	WHERE id = ?
	LIMIT 1`

	b, err := json.Marshal(banner)

	if err != nil {
		logger.WithField("location", "NewBanner").Error(err)

		return err
	}

	_, err = env.DB.Exec(query, string(b), id)

	if err != nil {
		logger.WithField("location", "NewBanner").Error(err)
		return err
	}

	return nil
}
