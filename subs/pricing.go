package subs

import (
	"strings"

	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/types/util"
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
func (p *Plan) Validate() *view.Reason {

	if p.NetPrice <= 0 {
		reason := view.NewReason()
		reason.Field = "netPrice"
		reason.Code = view.CodeInvalid
		reason.SetMessage("Net price must be greater than 0")

		return reason
	}

	if r := util.RequireNotEmptyWithMax(p.Description, 128, "description"); r != nil {
		return r
	}

	return util.OptionalMaxLen(p.Message, 128, "message")
}

// Pricing is an alias to a map of Plan.
type Pricing map[string]Plan

// Validate validates if pricing plans are valid.
func (p Pricing) Validate() *view.Reason {
	stdYear, ok := p[keyStdYear]

	if !ok {
		reason := view.NewReason()
		reason.Field = keyStdYear
		reason.Code = view.CodeMissingField

		return reason
	}

	if r := stdYear.Validate(); r != nil {
		r.Field = keyStdYear + "." + r.Field
		return r
	}

	if stdMonth, ok := p[keyStdMonth]; ok {
		if r := stdMonth.Validate(); r != nil {
			r.Field = keyStdMonth + "." + r.Field
			return r
		}
	}

	prmYear, ok := p[keyPrmYear]

	if !ok {
		reason := view.NewReason()
		reason.Field = keyPrmYear
		reason.Code = view.CodeMissingField

		return reason
	}

	if r := prmYear.Validate(); r != nil {
		r.Field = keyPrmYear + "." + r.Field
		return r
	}

	return nil
}
