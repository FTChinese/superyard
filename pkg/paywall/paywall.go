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
// Usually you have to retrieve the Product and Plans
// separately and assemble them.
type ProductExpanded struct {
	Product
	Plans []DiscountedPlan `json:"plans"`
}

func GroupPlans(plans []DiscountedPlan) map[string][]DiscountedPlan {
	var g = make(map[string][]DiscountedPlan)

	for _, v := range plans {
		if found, ok := g[v.ProductID]; ok {
			found = append(found, v)
		} else {
			g[v.ProductID] = []DiscountedPlan{v}
		}
	}

	return g
}

func BuildPaywallProducts(prods []Product, plans []DiscountedPlan) []ProductExpanded {
	groupedPlans := GroupPlans(plans)

	var result = make([]ProductExpanded, 0)

	for _, prod := range prods {
		gPlans, ok := groupedPlans[prod.ID]

		if !ok {
			gPlans = []DiscountedPlan{}
		}

		result = append(result, ProductExpanded{
			Product: prod,
			Plans:   gPlans,
		})
	}

	return result
}

type BannerInput struct {
	Heading    string      `json:"heading" db:"heading"`
	CoverURL   null.String `json:"coverUrl" db:"coverUrl"`
	SubHeading null.String `json:"subHeading" db:"sub_heading"`
	Content    null.String `json:"content" db:"content"`
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

type PromoInput struct {
	Banner
	Period
}

type Promo struct {
	ID string `json:"id" db:"promo_id"`
	Banner
	Period
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
}

// NewPromo create a new promotion based on input data.
func NewPromo(input PromoInput, creator string) Promo {
	return Promo{
		ID:         genPromoID(),
		Banner:     input.Banner,
		Period:     input.Period,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  creator,
	}
}