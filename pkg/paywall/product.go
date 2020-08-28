package paywall

import (
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
	CreatedBy  string      `json:"createdBy" db:"created_by"`
}

// Update modifies an existing product.
func (p Product) Update(input ProductInput) Product {
	p.Heading = input.Heading
	p.Description = input.Description
	p.SmallPrint = input.SmallPrint

	return p
}
