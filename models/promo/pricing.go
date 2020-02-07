package promo

import (
	"gitlab.com/ftchinese/superyard/models/validator"
	"strings"

	"github.com/FTChinese/go-rest/enum"
)

const (
	keyStdYear  = "standard_year"
	keyStdMonth = "standard_month"
	keyPrmYear  = "premium_year"
)

// Plan contains details of subs plan.
type Plan struct {
	Tier      enum.Tier  `json:"tier"`
	Cycle     enum.Cycle `json:"cycle"`
	ListPrice float64    `json:"listPrice"`
	NetPrice  float64    `json:"netPrice"`
	// For wxpay, this is used as `body` parameter;
	// For alipay, this is used as `subject` parameter.
	Description string `json:"description"` // required, max 128 chars
	// For wxpay, this is used as `detail` parameter;
	// For alipay, this is used as `body` parameter.
	Message string `json:"message"`
}

// Sanitize removes leading and trailing spaces of string fields.
func (p *Plan) Sanitize() {
	p.Description = strings.TrimSpace(p.Description)
	p.Message = strings.TrimSpace(p.Description)
}

// Validate validates if a plan is valid.
func (p *Plan) Validate() *validator.InputError {

	if p.NetPrice <= 0 {
		return &validator.InputError{
			Message: "Net price must be greater than 0",
			Field:   "netPrice",
			Code:    validator.CodeInvalid,
		}
	}

	ie := validator.New("description").Required().Max(128).Validate(p.Description)
	if ie != nil {
		return ie
	}

	return validator.New("message").Max(128).Validate(p.Message)
}

// Pricing is an alias to a map of Plan.
type Pricing map[string]Plan

// Validate validates if pricing plans are valid.
func (p Pricing) Validate() *validator.InputError {
	stdYear, ok := p[keyStdYear]

	if !ok {
		return &validator.InputError{
			Message: "Missing plan for yearly standard edition",
			Field:   keyStdYear,
			Code:    validator.CodeInvalid,
		}
	}

	if ie := stdYear.Validate(); ie != nil {
		ie.Field = keyStdYear + "." + ie.Field
		return ie
	}
	stdMonth, ok := p[keyStdMonth]

	if !ok {
		return &validator.InputError{
			Message: "Missing plan for monthly standard edition",
			Field:   keyStdMonth,
			Code:    validator.CodeInvalid,
		}
	}

	if ie := stdMonth.Validate(); ie != nil {
		ie.Field = keyStdMonth + "." + ie.Field
		return ie
	}

	prmYear, ok := p[keyPrmYear]

	if !ok {
		return &validator.InputError{
			Message: "Missing plan for yearly premium edition",
			Field:   keyPrmYear,
			Code:    validator.CodeInvalid,
		}
	}

	if ie := prmYear.Validate(); ie != nil {
		ie.Field = keyPrmYear + "." + ie.Field
		return ie
	}

	return nil
}
