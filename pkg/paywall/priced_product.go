package paywall

import (
	"encoding/json"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
)

// PricedProductInput defines the input data to create a product,
// with optional plans.
type PricedProductInput struct {
	ProductInput
	// Plans created this way have only price, cycle, description fields. Tier is dependent on Product.Tier
	Plans []PlanInput `json:"plans"`
}

// PricedProduct is a product containing pricing plans.
// The plan does not contain discount.
// It is used in two places:
// * Show a list of product. Each item has an overview of how many plans but won't show the details of those plans.
// * Creating a new product. Opitional prices can be created.
type PricedProduct struct {
	Product
	Plans []Plan `json:"plans"`
}

// NewPricedProduct creates a new product instance based on input.
func NewPricedProduct(input PricedProductInput, creator string) PricedProduct {

	prodID := GenProductID()

	var plans = make([]Plan, 0)
	for _, v := range input.Plans {
		// Don't forget to add product id to plan.
		// Call NewPlan() won't add it since it assumes to
		// be provided by client.
		v.ProductID = prodID
		v.Tier = input.Tier
		plans = append(plans, NewPlan(v, creator))
	}

	return PricedProduct{
		Product: Product{
			ID:           prodID,
			ProductInput: input.ProductInput,
			IsActive:     false,
			CreatedUTC:   chrono.TimeNow(),
			UpdatedUTC:   chrono.TimeNow(),
			CreatedBy:    creator,
		},
		Plans: plans,
	}
}

// Validate checks whether the request to product with
// optional plans are valid.
func (p PricedProduct) Validate() *render.ValidationError {
	ve := p.ProductInput.Validate()
	if ve != nil {
		return ve
	}

	for _, v := range p.Plans {
		if ve := v.PlanInput.Validate(); ve != nil {
			return ve
		}
	}

	return nil
}

// PricedProductSchema is used to hold db scan data for a list of product with plans retrieved as a JSON string.
type PricedProductSchema struct {
	Product
	Plans null.String `db:"plans"`
}

// PricedProduct converts data retrieve from db.
func (s PricedProductSchema) PricedProduct() (PricedProduct, error) {
	var plans = make([]Plan, 0)

	err := json.Unmarshal([]byte(s.Plans.String), &plans)
	if err != nil {
		return PricedProduct{}, err
	}

	return PricedProduct{
		Product: s.Product,
		Plans:   plans,
	}, nil
}
