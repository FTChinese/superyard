package products

import (
	"database/sql"
	"github.com/FTChinese/superyard/pkg/paywall"
)

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

func (env Env) ListPlansOfProduct(id string) ([]paywall.DiscountedPlan, error) {
	schemas := make([]paywall.DiscountedPlanSchema, 0)
	dPlans := make([]paywall.DiscountedPlan, 0)

	err := env.db.Select(&schemas, paywall.StmtPlansOfProduct, id)

	if err != nil {
		return dPlans, err
	}

	for _, v := range schemas {
		dPlans = append(dPlans, v.DiscountedPlan())
	}

	return dPlans, nil
}
