package products

import (
	"database/sql"
	"github.com/FTChinese/superyard/pkg/paywall"
)

type bannerResult struct {
	value paywall.Banner
	err   error
}

func (env Env) asyncLoadBanner(id int64) <-chan bannerResult {
	c := make(chan bannerResult)

	go func() {
		b, err := env.LoadBanner(id)
		c <- bannerResult{
			value: b,
			err:   err,
		}
	}()

	return c
}

// retrievePaywallPromo selects a promo whose id is set on the specified banner.
func (env Env) retrievePaywallPromo(bannerID int64) (paywall.Promo, error) {
	var promo paywall.Promo
	err := env.db.Get(&promo, paywall.StmtPaywallPromo, bannerID)
	if err != nil && err != sql.ErrNoRows {
		return paywall.Promo{}, err
	}

	return promo, nil
}

type promoResult struct {
	value paywall.Promo
	err   error
}

func (env Env) asyncLoadPaywallPromo(bannerID int64) <-chan promoResult {
	c := make(chan promoResult)

	go func() {
		defer close(c)
		p, err := env.retrievePaywallPromo(bannerID)

		c <- promoResult{
			value: p,
			err:   err,
		}
	}()

	return c
}

// retrievePaywallProducts retrieves all products present on paywall.
// Those products does not include its pricing plans.
// You need to zip them with the result from retrievePaywallPlans.
func (env Env) retrievePaywallProducts() ([]paywall.Product, error) {
	var products = make([]paywall.Product, 0)

	err := env.db.Select(&products, paywall.StmtPaywallProducts)
	if err != nil {
		return nil, err
	}

	return products, nil
}

type productsResult struct {
	value []paywall.Product
	err   error
}

func (env Env) asyncLoadPaywallProducts() <-chan productsResult {
	c := make(chan productsResult)

	go func() {
		p, err := env.retrievePaywallProducts()

		c <- productsResult{
			value: p,
			err:   err,
		}
	}()

	return c
}

// retrievePaywallPlans retrieves all plans of the specified productIDs array.
// The productIDs should be the ids retrieved by retrievePaywallProducts
func (env Env) retrievePaywallPlans() ([]paywall.ExpandedPlan, error) {
	schemas := make([]paywall.ExpandedPlanSchema, 0)
	var plans = make([]paywall.ExpandedPlan, 0)

	err := env.db.Select(&schemas, paywall.StmtPaywallPlans)

	if err != nil {
		return nil, err
	}

	for _, v := range schemas {
		plans = append(plans, v.ExpandedPlan())
	}

	return plans, nil
}

type plansResult struct {
	value []paywall.ExpandedPlan
	err   error
}

func (env Env) asyncLoadPaywallPlans() <-chan plansResult {
	c := make(chan plansResult)

	go func() {
		p, err := env.retrievePaywallPlans()

		c <- plansResult{
			value: p,
			err:   err,
		}
	}()

	return c
}

func (env Env) LoadPaywall(id int64) (paywall.Paywall, error) {
	bannerCh, promoCh, productsCh, plansCh := env.asyncLoadBanner(id), env.asyncLoadPaywallPromo(id), env.asyncLoadPaywallProducts(), env.asyncLoadPaywallPlans()

	bannerRes, promoRes, productsRes, plansRes := <-bannerCh, <-promoCh, <-productsCh, <-plansCh

	if bannerRes.err != nil {
		return paywall.Paywall{}, bannerRes.err
	}

	if promoRes.err != nil {
		return paywall.Paywall{}, promoRes.err
	}

	if productsRes.err != nil {
		return paywall.Paywall{}, productsRes.err
	}

	if plansRes.err != nil {
		return paywall.Paywall{}, plansRes.err
	}

	return paywall.Paywall{
		Banner: bannerRes.value,
		Promo:  promoRes.value,
		Products: paywall.BuildPaywallProducts(
			productsRes.value,
			plansRes.value,
		),
	}, nil
}
