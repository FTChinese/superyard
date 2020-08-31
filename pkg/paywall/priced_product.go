package paywall

import (
	"fmt"
	"github.com/FTChinese/go-rest/render"
)

// PricedProductInput defines the input data to create a product, with optional plans.
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

	product := NewProduct(input.ProductInput, creator)

	var plans = make([]Plan, 0)
	for _, v := range input.Plans {
		// Don't forget to add product id to plan.
		// Call NewPlan() won't add it since it assumes to
		// be provided by client.
		v.ProductID = product.ID
		plans = append(plans, product.NewPlan(v, creator))
	}

	return PricedProduct{
		Product: product,
		Plans:   plans,
	}
}

// Validate checks whether the request to create product with optional plans are valid.
func (p PricedProduct) Validate() *render.ValidationError {
	ve := p.ProductInput.Validate()
	if ve != nil {
		return ve
	}

	for i, v := range p.Plans {
		if ve := v.PlanInput.Validate(); ve != nil {
			// Here we modified the validation error's field to the path of error field.
			// `plans` is the top level field json tag;
			// `i` is the position of the array;
			// and finally it is the field name errored.
			ve.Field = fmt.Sprintf("%s.%d.%s", "plans", i, ve.Field)
			return ve
		}
		if ve := v.IsCycleMismatched(); ve != nil {
			return ve
		}
	}

	return nil
}

// ListedProduct is used to show a list of products with each product's total plans count.
type ListedProduct struct {
	Product
	PlanCount int64 `json:"planCount" db:"plan_count"`
}
