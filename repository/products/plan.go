package products

import "github.com/FTChinese/superyard/pkg/paywall"

func (env Env) CreatePlan(p paywall.Plan) error {
	_, err := env.db.NamedExec(paywall.StmtCreatePlan, p)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadPlan(id string) (paywall.Plan, error) {
	var plan paywall.Plan

	err := env.db.Get(&plan, paywall.StmtPlan, id)

	if err != nil {
		return plan, err
	}

	return plan, nil
}

func (env Env) ActivatePlan(plan paywall.Plan) error {
	_, err := env.db.NamedExec(paywall.StmtActivatePlan, plan)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) ListPlansOfProduct(id string) ([]paywall.DiscountedPlan, error) {
	schemas := make([]paywall.DiscountedPlanSchema, 0)
	dPlans := make([]paywall.DiscountedPlan, 0)

	err := env.db.Select(&schemas, paywall.StmtPlansOfProduct)

	if err != nil {
		return dPlans, err
	}

	for _, v := range schemas {
		dPlans = append(dPlans, v.DiscountedPlan())
	}

	return dPlans, nil
}
