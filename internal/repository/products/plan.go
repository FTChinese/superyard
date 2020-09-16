package products

import (
	"database/sql"
	"github.com/FTChinese/superyard/pkg/paywall"
)

// CreatePlan saves a new plan for a product.
func (env Env) CreatePlan(p paywall.Plan) error {
	_, err := env.db.NamedExec(paywall.StmtCreatePlan, p)

	if err != nil {
		return err
	}

	return nil
}

// LoadPlan retrieves a plan. It contains only columns from the plan table.
// Not discount info is included.
func (env Env) LoadPlan(id string) (paywall.Plan, error) {
	var plan paywall.Plan

	err := env.db.Get(&plan, paywall.StmtPlan, id)

	if err != nil {
		return plan, err
	}

	return plan, nil
}

// ActivatePlan sets a plan as active under a product.
// By active we mean this plan will be presented on paywall
// when this product is put on paywall.
// A product could have as many plans linked as you like,
// but only two (for product tier standard) or one (for premium)
// set as active.
func (env Env) ActivatePlan(plan paywall.Plan) error {
	_, err := env.db.NamedExec(paywall.StmtActivatePlan, plan)
	if err != nil {
		return err
	}

	return nil
}

// ProductHasActivePlan checks if a product has any active plans set for it.
// If a product does not have any plans, we should disallow
// it being put on paywall.
func (env Env) ProductHasActivePlan(productID string) (bool, error) {
	var ok bool
	err := env.db.Get(&ok, paywall.StmtHasActivePlan, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return ok, nil
}

// ListPlansOfProduct retrieves all plans under a product.
// Each plans has discount attached to it.
func (env Env) ListPlansOfProduct(productID string) ([]paywall.ExpandedPlan, error) {
	schemas := make([]paywall.ExpandedPlanSchema, 0)
	dPlans := make([]paywall.ExpandedPlan, 0)

	err := env.db.Select(&schemas, paywall.StmtPlansOfProduct, productID)

	if err != nil {
		return dPlans, err
	}

	for _, v := range schemas {
		dPlans = append(dPlans, v.ExpandedPlan())
	}

	return dPlans, nil
}

func (env Env) ListPlansOnPaywall() ([]paywall.Plan, error) {
	plans := make([]paywall.Plan, 0)

	err := env.db.Select(&plans, paywall.StmtListPlansOnPaywall)
	if err != nil {
		return nil, err
	}

	return plans, nil
}

func (env Env) PaywallPlanByEdition(edition paywall.Edition) (paywall.Plan, error) {
	var plan paywall.Plan

	err := env.db.Get(&plan, paywall.StmtPaywallPlan, edition.Tier, edition.Cycle)
	if err != nil {
		return paywall.Plan{}, err
	}

	return plan, nil
}