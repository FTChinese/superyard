package paywall

import (
	"encoding/json"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strings"
)

// ProductInput defines the request data to create a new
// product.
type ProductInput struct {
	Tier        enum.Tier   `json:"tier" db:"tier"`
	Heading     string      `json:"heading" db:"heading"`
	Description null.String `json:"description" db:"description"`
	SmallPrint  null.String `json:"smallPrint" db:"small_print"`
}

// Validate checks fields to create or update a product.
func (p *ProductInput) Validate() *render.ValidationError {
	p.Heading = strings.TrimSpace(p.Heading)
	p.Description.String = strings.TrimSpace(p.Description.String)
	p.SmallPrint.String = strings.TrimSpace(p.SmallPrint.String)

	if p.Tier == enum.TierNull {
		return &render.ValidationError{
			Message: "Tier could not be null",
			Field:   "tier",
			Code:    render.CodeMissingField,
		}
	}
	return validator.
		New("heading").
		Required().
		Validate(p.Heading)
}

// Product defines a product without plans.
type Product struct {
	ID string `json:"id" db:"product_id"`
	ProductInput
	IsActive   bool        `json:"isActive" db:"is_active"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
	CreatedBy  string      `json:"createdByy" db:"created_by"`
}

// Update modifies an existing product.
func (p Product) Update(input ProductInput) Product {
	p.Heading = input.Heading
	p.Description = input.Description
	p.SmallPrint = input.SmallPrint

	return p
}

// PricedProductInput defines the input data to create a product,
// with optional plans.
type PricedProductInput struct {
	ProductInput
	// Plans created this way have only price, cycle, description fields. Tier is dependent on Product.Tier
	Plans []PlanInput `json:"plans"`
}

// PricedProduct is a product containing pricing plans.
// The plan does not contain discount.
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
