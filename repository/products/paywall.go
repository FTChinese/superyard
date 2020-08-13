package products

import (
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/guregu/null"
	"strings"
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
		return p, err
	}

	return p, nil
}

func (env Env) ListActiveProducts() ([]paywall.Product, error) {
	var products = make([]paywall.Product, 0)

	err := env.db.Select(&products, paywall.StmtActiveProducts)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (env Env) ListActivePlans(productIDs []string) ([]paywall.DiscountedPlan, error) {
	var plans = make([]paywall.DiscountedPlan, 0)
	idSet := strings.Join(productIDs, ",")

	err := env.db.Select(&plans, paywall.StmtActivePlans, idSet)

	if err != nil {
		return nil, err
	}

	return plans, nil
}

func (env Env) LoadPaywallProducts() ([]paywall.ProductExpanded, error) {

	prods, err := env.ListActiveProducts()
	if err != nil {
		return nil, err
	}

	if len(prods) == 0 {
		return []paywall.ProductExpanded{}, nil
	}

	prodIds := make([]string, 0)
	for _, v := range prods {
		prodIds = append(prodIds, v.ID)
	}

	plans, err := env.ListActivePlans(prodIds)
	if err != nil {
		return nil, err
	}

	return paywall.BuildPaywallProducts(prods, plans), nil
}
