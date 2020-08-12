package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
)

type DiscountInput struct {
	PriceOff null.Int `json:"priceOff" db:"price_off"`
	Percent  null.Int `json:"percent" db:"percent"`
	Period
}

type Discount struct {
	ID     null.String `json:"id" db:"discount_id"`
	PlanID null.String `json:"planId" db:"plan_id"` // This is used only to save data.
	DiscountInput
}

func NewDiscount(input DiscountInput, planID string) Discount {
	return Discount{
		ID:            null.StringFrom(genDiscountID()),
		PlanID:        null.StringFrom(planID),
		DiscountInput: input,
	}
}

// Validate checks whether request data to create a discount
// are valid.
// Required fields:
// PriceOff
// StartUTC
// EndUTC
func (d Discount) Validate() *render.ValidationError {

	if d.PriceOff.IsZero() || d.PriceOff.Int64 <= 0 {
		return &render.ValidationError{
			Message: "priceOff is required",
			Field:   "priceOff",
			Code:    render.CodeMissingField,
		}
	}

	if d.StartUTC.IsZero() {
		return &render.ValidationError{
			Message: "startUtc is required",
			Field:   "startUtc",
			Code:    render.CodeMissingField,
		}
	}

	if d.EndUTC.IsZero() {
		return &render.ValidationError{
			Message: "endUtc is required",
			Field:   "endUtc",
			Code:    render.CodeMissingField,
		}
	}

	if d.StartUTC.After(d.EndUTC.Time) {
		return &render.ValidationError{
			Message: "start time must be earlier than end time",
			Field:   "startUtc",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

// DiscountSchema is used to insert a discount row.
type DiscountSchema struct {
	Discount
	CreatedUTC chrono.Time `db:"created_utc"`
	CreatedBy  string      `db:"created_by"`
}

func NewDiscountSchema(d Discount, creator string) DiscountSchema {
	return DiscountSchema{
		Discount:   d,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  creator,
	}
}
