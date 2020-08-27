package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strings"
)

// BannerInput is used to create a new banner.
type BannerInput struct {
	Heading    string      `json:"heading" db:"heading"`
	SubHeading null.String `json:"subHeading" db:"sub_heading"`
	CoverURL   null.String `json:"coverUrl" db:"cover_url"`
	Content    null.String `json:"content" db:"content"`
}

func (b BannerInput) Validate() *render.ValidationError {
	b.Heading = strings.TrimSpace(b.Heading)
	b.CoverURL.String = strings.TrimSpace(b.CoverURL.String)
	b.SubHeading.String = strings.TrimSpace(b.SubHeading.String)
	b.Content.String = strings.TrimSpace(b.Content.String)

	return validator.New("heading").Required().Validate(b.Heading)
}

// Banner is the banner data shown on paywall.
type Banner struct {
	ID int64 `json:"id" db:"banner_id"`
	BannerInput
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
	PromoID    null.String `json:"promoId" db:"promo_id"`
}

// NewBanner creates a new Banner instance based on input data.
func NewBanner(input BannerInput, creator string) Banner {
	return Banner{
		BannerInput: input,
		CreatedUTC:  chrono.TimeNow(),
		UpdatedUTC:  chrono.TimeNow(),
		CreatedBy:   creator,
		PromoID:     null.String{},
	}
}

func (b Banner) Update(input BannerInput) Banner {
	b.Heading = input.Heading
	b.CoverURL = input.CoverURL
	b.SubHeading = input.SubHeading
	b.Content = input.Content

	return b
}
