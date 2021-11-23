package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strings"
)

// PlanInput represents the data used to create a new plan.
// A new plan is always created under a certain product.
// Therefore the input data does not have tier field.
type PlanInput struct {
	Cycle       enum.Cycle  `json:"cycle" db:"cycle"`
	Description null.String `json:"description" db:"description"`
	Price       float64     `json:"price" db:"price"`
	ProductID   string      `json:"productId" db:"product_id"`
}

// Validate checks whether the input data to create a new plan is valid.
// `productTier` is used to specify for which edition of product this plan is created.
// Premium product is not allowed to have a monthly pricing plan.
func (p *PlanInput) Validate() *render.ValidationError {

	p.Description.String = strings.TrimSpace(p.Description.String)

	ve := validator.New("productId").Required().Validate(p.ProductID)
	if ve != nil {
		return ve
	}

	if p.Price <= 0 {
		return &render.ValidationError{
			Message: "Price could not be below 0",
			Field:   "price",
			Code:    render.CodeInvalid,
		}
	}

	if p.Cycle == enum.CycleNull {
		return &render.ValidationError{
			Message: "Invalid cycle",
			Field:   "cycle",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

// Plan do not contain the discount data.
type Plan struct {
	ID string `json:"id" db:"plan_id"`
	PlanInput
	Tier       enum.Tier   `json:"tier" db:"tier"`
	IsActive   bool        `json:"isActive" db:"is_active"`
	CreatedUTC chrono.Time `json:"created_utc" db:"created_utc"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
}

// IsCycleMismatched checks whether this plan's
// billing cycle is applicable to its tier.
func (p Plan) IsCycleMismatched() *render.ValidationError {
	if p.Tier == enum.TierPremium && p.Cycle == enum.CycleMonth {
		return &render.ValidationError{
			Message: "Billing cycle month is not applicable to premium plan",
			Field:   "cycle",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

// ExpandedPlan is used to output a plan with optional discount.
type ExpandedPlan struct {
	Plan
	Discount Discount `json:"discount"`
}

// ExpandedPlanSchema is used to retrieve a plan with discount.
type ExpandedPlanSchema struct {
	Plan
	Discount
}

// ExpandedPlan turns the retrieved data to ExpandedPlan.
func (s ExpandedPlanSchema) ExpandedPlan() ExpandedPlan {
	return ExpandedPlan{
		Plan:     s.Plan,
		Discount: s.Discount,
	}
}
