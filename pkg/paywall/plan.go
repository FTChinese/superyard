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
type PlanInput struct {
	ProductID   string      `json:"productId" db:"product_id"`
	Price       float64     `json:"price" db:"price"`
	Tier        enum.Tier   `json:"tier" db:"tier"`
	Cycle       enum.Cycle  `json:"cycle" db:"cycle"`
	Description null.String `json:"description" db:"description"`
}

// Validate checks whether the input data to create a new plan is valid.
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

	if p.Tier == enum.TierNull {
		return &render.ValidationError{
			Message: "Plan tier is not valid",
			Field:   "tier",
			Code:    render.CodeInvalid,
		}
	}

	if p.Cycle == enum.CycleNull {
		return &render.ValidationError{
			Message: "Plan cycle is not valid",
			Field:   "cycle",
			Code:    render.CodeInvalid,
		}
	}

	if p.Tier == enum.TierPremium && p.Cycle == enum.CycleMonth {
		return &render.ValidationError{
			Message: "Billing cycle month is not applicable to premium plan",
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
	IsActive   bool        `json:"isActive" db:"is_active"`
	CreatedUTC chrono.Time `json:"created_utc" db:"created_utc"`
	CreatedBy  string      `json:"createdBy" db:"created_by"`
}

// NewPlan creates a new Plan from input data.
func NewPlan(input PlanInput, creator string) Plan {
	return Plan{
		ID:         genPlanID(),
		PlanInput:  input,
		IsActive:   false,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  creator,
	}
}

// ExpandedPlan is used to output a plan with optional discount.
type ExpandedPlan struct {
	Plan
	Discount Discount `json:"discount"`
}

// GroupPlans groups plans into a map by product id.
func GroupPlans(plans []ExpandedPlan) map[string][]ExpandedPlan {
	var g = make(map[string][]ExpandedPlan)

	for _, v := range plans {
		found, ok := g[v.ProductID]
		if ok {
			found = append(found, v)
		} else {
			found = []ExpandedPlan{v}
		}
		g[v.ProductID] = found
	}

	return g
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
