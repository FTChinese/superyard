package products

import (
	"github.com/FTChinese/superyard/pkg/paywall"
)

func (env Env) CreatePricedProduct(p paywall.PricedProduct) error {

	tx, err := env.dbs.Write.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(paywall.StmtCreateProduct, p.Product)

	if err != nil {
		return err
	}

	if len(p.Plans) > 0 {

		for _, v := range p.Plans {
			_, err := tx.NamedExec(paywall.StmtCreatePlan, v)
			if err != nil {
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (env Env) LoadProduct(id string) (paywall.Product, error) {
	var p paywall.Product

	err := env.dbs.Read.Get(&p, paywall.StmtProduct, id)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (env Env) UpdateProduct(prod paywall.Product) error {
	_, err := env.dbs.Write.NamedExec(paywall.StmtUpdateProduct, prod)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) ActivateProduct(prod paywall.Product) error {
	_, err := env.dbs.Write.NamedExec(paywall.StmtActivateProduct, prod)

	if err != nil {
		return err
	}

	return nil
}

// ListProducts list all products, with each product's plans attached without discount.
// A product's plans are retrieve using JSON_ARRAYAGG.
func (env Env) ListProducts() ([]paywall.ListedProduct, error) {
	products := make([]paywall.ListedProduct, 0)

	err := env.dbs.Read.Select(&products, paywall.StmtListPricedProducts)

	if err != nil {
		return products, err
	}

	return products, nil
}
