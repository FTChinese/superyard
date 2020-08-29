package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strings"
)

type PromoInput struct {
	Heading    null.String `json:"heading" db:"heading"` // Required for input.
	SubHeading null.String `json:"subHeading" db:"sub_heading"`
	CoverURL   null.String `json:"coverUrl" db:"cover_url"`
	Content    null.String `json:"content" db:"content"`
	Terms      null.String `json:"terms" db:"terms_conditions"`
	Period
}

func (i *PromoInput) Validate() *render.ValidationError {
	i.Heading.String = strings.TrimSpace(i.Heading.String)
	i.CoverURL.String = strings.TrimSpace(i.CoverURL.String)
	i.SubHeading.String = strings.TrimSpace(i.SubHeading.String)
	i.Content.String = strings.TrimSpace(i.Content.String)
	i.Terms.String = strings.TrimSpace(i.Terms.String)

	if !i.EndUTC.After(i.StartUTC.Time) {
		return &render.ValidationError{
			Message: "End time should after start time",
			Field:   "endUtc",
			Code:    render.CodeInvalid,
		}
	}

	return validator.New("heading").Required().Validate(i.Heading.String)
}

type Promo struct {
	ID null.String `json:"id" db:"promo_id"`
	PromoInput
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	CreatedBy  null.String `json:"createdBy" db:"created_by"`
}

// NewPromo create a new promotion based on input data.
func NewPromo(input PromoInput, creator string) Promo {
	return Promo{
		ID:         null.StringFrom(genPromoID()),
		PromoInput: input,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  null.StringFrom(creator),
	}
}
