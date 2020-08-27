package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

type PromoInput struct {
	BannerInput
	Terms null.String `json:"terms" db:"term_conditions"`
	Period
}

type Promo struct {
	ID string `json:"id" db:"promo_id"`
	BannerInput
	Period
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
}

// NewPromo create a new promotion based on input data.
func NewPromo(input PromoInput, creator string) Promo {
	return Promo{
		ID:          genPromoID(),
		BannerInput: input.BannerInput,
		Period:      input.Period,
		CreatedUTC:  chrono.TimeNow(),
		CreatedBy:   creator,
	}
}
