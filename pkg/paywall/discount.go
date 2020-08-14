package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
)

type DiscountInput struct {
	PriceOff null.Float `json:"priceOff" db:"price_off"`
	Percent  null.Int   `json:"percent" db:"percent"`
	Period
}

// Validate checks whether request data to create a discount
// are valid.
func (d DiscountInput) Validate() *render.ValidationError {

	if d.PriceOff.IsZero() || d.PriceOff.Float64 <= 0 {
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

// Discount is a plan's discount. User for output.
type Discount struct {
	// The id fields started with Disc to avoid conflict when used in DiscountedPlanSchema.
	DiscID     null.String `json:"id" db:"discount_id"`
	DiscPlanID null.String `json:"planId" db:"discounted_plan_id"`
	DiscountInput
}

// DiscountSchema is used to insert a discount row.
// Used in db insert. Every row is immutable.
type DiscountSchema struct {
	Discount
	CreatedUTC chrono.Time `db:"created_utc"`
	CreatedBy  string      `db:"created_by"`
}

func NewDiscountSchema(input DiscountInput, planID, creator string) DiscountSchema {
	return DiscountSchema{
		Discount: Discount{
			DiscID:        null.StringFrom(genDiscountID()),
			DiscPlanID:    null.StringFrom(planID),
			DiscountInput: input,
		},
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  creator,
	}
}
