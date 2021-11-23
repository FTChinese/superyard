package products

import (
	"database/sql"
	"github.com/FTChinese/superyard/pkg/paywall"
)

// ProductHasActivePlan checks if a product has any active plans set for it.
// If a product does not have any plans, we should disallow
// it being put on paywall.
func (env Env) ProductHasActivePlan(productID string) (bool, error) {
	var ok bool
	err := env.dbs.Read.Get(&ok, paywall.StmtHasActivePlan, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return ok, nil
}

func (env Env) ListPlansOnPaywall() ([]paywall.Plan, error) {
	plans := make([]paywall.Plan, 0)

	err := env.dbs.Read.Select(&plans, paywall.StmtListPlansOnPaywall)
	if err != nil {
		return nil, err
	}

	return plans, nil
}

func (env Env) PaywallPlanByEdition(edition paywall.Edition) (paywall.Plan, error) {
	var plan paywall.Plan

	err := env.dbs.Read.Get(&plan, paywall.StmtPaywallPlan, edition.Tier, edition.Cycle)
	if err != nil {
		return paywall.Plan{}, err
	}

	return plan, nil
}
