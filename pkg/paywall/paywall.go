package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/validator"
	"strings"
)

// ProductExpanded defines the complete data of a product.
// This is mostly used when to compose the paywall data
// However it not easy to retrieve all its data in one shot.
// Usually you have to retrieve the BaseProduct and Plans
// separately and assemble them.
type ProductExpanded struct {
	BaseProduct
	Plans []Plan `json:"plans"`
}

type BannerInput struct {
	Heading    string      `json:"heading" db:"heading"`
	CoverURL   null.String `json:"coverUrl" db:"coverUrl"`
	SubHeading null.String `json:"subHeading" db:"sub_heading"`
	Content    null.String `json:"content" db:"content"`
}

func (b BannerInput) NewBanner(creator string) Banner {
	return Banner{
		BannerInput: b,
		CreatedUTC:  chrono.TimeNow(),
		UpdatedUTC:  chrono.TimeNow(),
		CreatedBy:   creator,
		PromoID:     null.String{},
	}
}

func (b BannerInput) Validate() *render.ValidationError {
	b.Heading = strings.TrimSpace(b.Heading)
	b.CoverURL.String = strings.TrimSpace(b.CoverURL.String)
	b.SubHeading.String = strings.TrimSpace(b.SubHeading.String)
	b.Content.String = strings.TrimSpace(b.Content.String)

	return validator.New("heading").Required().Validate(b.Heading)
}

type Banner struct {
	ID int64 `json:"id" db:"id"`
	BannerInput
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
	PromoID    null.String `json:"promoId" db:"promo_id"`
}

func (b Banner) Update(input BannerInput) Banner {
	b.Heading = input.Heading
	b.CoverURL = input.CoverURL
	b.SubHeading = input.SubHeading
	b.Content = input.Content

	return b
}

type PromoInput struct {
	Banner
	Period
}

func (p PromoInput) NewPromo(creator string) Promo {
	return Promo{
		ID:         genPromoID(),
		Banner:     p.Banner,
		Period:     p.Period,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  creator,
	}
}

type Promo struct {
	ID string `json:"id" db:"promo_id"`
	Banner
	Period
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
}
