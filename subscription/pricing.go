package subscription

import (
	"encoding/json"
	"strings"

	"gitlab.com/ftchinese/backyard-api/util"
)

var tiers = map[string]int{
	"standard": 0,
	"premium":  1,
}

var cycles = map[string]int{
	"year":  0,
	"month": 1,
}

const (
	keyStdYear  = "standard_year"
	keyStdMonth = "standard_month"
	keyPrmYear  = "premium_year"
)

// PromoPlan contains details of subscription plan.
type PromoPlan struct {
	Tier  string  `json:"tier"`
	Cycle string  `json:"cycle"`
	Price float64 `json:"price"`
	ID    int
	// For wxpay, this is used as `body` parameter;
	// For alipay, this is used as `subject` parameter.
	Description string `json:"description"` // required, max 128 chars
	// For wxpay, this is used as `detail` parameter;
	// For alipay, this is used as `body` parameter.
	Message string `json:"message"`
	Ignore  bool   `json:"ignore,omitempty"`
}

// Sanitize removes leading and trailing spaces of string fields.
func (p *PromoPlan) Sanitize() {
	p.Tier = strings.TrimSpace(p.Tier)
	p.Cycle = strings.TrimSpace(p.Cycle)
	p.Description = strings.TrimSpace(p.Description)
	p.Message = strings.TrimSpace(p.Description)
}

// Validate validates if a plan is valid.
func (p *PromoPlan) Validate() *util.Reason {
	if r := util.RequireNotEmpty(p.Tier, "tier"); r != nil {
		return r
	}

	if r := util.RequireNotEmpty(p.Cycle, "cycle"); r != nil {
		return r
	}

	if _, ok := tiers[p.Tier]; !ok {
		reason := util.NewReason()
		reason.Field = "tier"
		reason.Code = util.CodeInvalid
		reason.SetMessage("Tier must be one of standard or premium")

		return reason
	}

	if _, ok := cycles[p.Cycle]; !ok {
		reason := util.NewReason()
		reason.Field = "cycle"
		reason.Code = util.CodeInvalid
		reason.SetMessage("Cycle must be one of year or month")
		return reason
	}

	if p.Price <= 0 {
		reason := util.NewReason()
		reason.Field = "price"
		reason.Code = util.CodeInvalid
		reason.SetMessage("Price must be greated than 0")

		return reason
	}

	if r := util.RequireNotEmptyWithMax(p.Description, 128, "description"); r != nil {
		return r
	}

	return util.OptionalMaxLen(p.Message, 128, "message")
}

// PromoPricing is an alias to a map of Plan.
type PromoPricing map[string]PromoPlan

// Validate validates if pricing plans are valid.
func (p PromoPricing) Validate() *util.Reason {
	stdYear, ok := p[keyStdYear]

	if !ok {
		reason := util.NewReason()
		reason.Field = keyStdYear
		reason.Code = util.CodeMissingField

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
		reason := util.NewReason()
		reason.Field = keyPrmYear
		reason.Code = util.CodeMissingField

		return reason
	}

	if r := prmYear.Validate(); r != nil {
		r.Field = keyPrmYear + "." + r.Field
		return r
	}

	return nil
}

// SavePricing set the pricing plans of a promotion schedule.
func (env Env) SavePricing(id int64, plans PromoPricing) error {
	query := `
	UPDATE premium.promotion_schedule
	SET plans = ?
	WHERE id = ?
	LIMIT 1`

	p, err := json.Marshal(plans)

	if err != nil {
		logger.WithField("location", "NewPricing").Error(err)
		return err
	}

	_, err = env.DB.Exec(query, string(p), id)

	if err != nil {
		logger.WithField("location", "NewPricing").Error(err)
		return err
	}

	return nil
}
