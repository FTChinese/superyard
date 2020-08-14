package products

import (
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/guregu/null"
)

func (env Env) CreateBanner(b paywall.Banner) error {
	_, err := env.db.NamedExec(paywall.StmtCreateBanner, b)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadBanner(id int64) (paywall.Banner, error) {
	var b paywall.Banner
	err := env.db.Get(&b, paywall.StmtBanner, id)
	if err != nil {
		return b, err
	}

	return b, err
}

func (env Env) UpdateBanner(b paywall.Banner) error {
	_, err := env.db.NamedExec(paywall.StmtUpdateBanner, b)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) DropBannerPromo(bannerID int64) error {
	_, err := env.db.Exec(paywall.StmtDropPromo, bannerID)
	if err != nil {
		return err
	}

	return nil
}

// CreatePromo creates a new promotion and apply it to banner immediately.
func (env Env) CreatePromo(bannerID int64, p paywall.Promo) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(paywall.StmtCreatePromo, p)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(paywall.StmtApplyPromo, paywall.Banner{
		ID:      bannerID,
		PromoID: null.StringFrom(p.ID),
	})
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (env Env) LoadPromo(id string) (paywall.Promo, error) {
	var p paywall.Promo

	err := env.db.Get(&p, paywall.StmtPromo, id)
	if err != nil {
		getLogger("LoadPromo").Error(err)
		return p, err
	}

	return p, nil
}

// listPaywallProducts retrieves all products present on paywall.
// Those products does not include its pricing plans.
// You need to zip them with the result from listPaywallPlans.
func (env Env) listPaywallProducts() ([]paywall.Product, error) {
	var products = make([]paywall.Product, 0)

	err := env.db.Select(&products, paywall.StmtPaywallProducts)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// listPaywallPlans retrieves all plans of the specified productIDs array.
// The productIDs should be the ids retrieved by listPaywallProducts
func (env Env) listPaywallPlans() ([]paywall.DiscountedPlan, error) {
	var plans = make([]paywall.DiscountedPlan, 0)

	err := env.db.Select(&plans, paywall.StmtPaywallPlans)

	if err != nil {
		return nil, err
	}

	return plans, nil
}

func (env Env) LoadPaywallProducts() ([]paywall.ProductExpanded, error) {

	prods, err := env.listPaywallProducts()
	if err != nil {
		return nil, err
	}

	if len(prods) == 0 {
		return []paywall.ProductExpanded{}, nil
	}

	plans, err := env.listPaywallPlans()
	if err != nil {
		return nil, err
	}

	return paywall.BuildPaywallProducts(prods, plans), nil
}
