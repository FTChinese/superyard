package products

import (
	"github.com/FTChinese/superyard/pkg/paywall"
)

func (env Env) CreatePricedProduct(p paywall.PricedProduct) error {

	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	logger.Info("Creating product...")
	_, err = tx.NamedExec(paywall.StmtCreateProduct, p.Product)

	if err != nil {
		getLogger("CreatePricedProduct").Error(err)
		return err
	}

	if len(p.Plans) > 0 {
		logger.Infof("Creating %d plans", len(p.Plans))

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

	err := env.db.Get(&p, paywall.StmtProduct, id)

	if err != nil {
		getLogger("LoadBaseProduct").Error(err)
		return p, err
	}

	return p, nil
}

func (env Env) UpdateProduct(prod paywall.Product) error {
	_, err := env.db.NamedExec(paywall.StmtUpdateProduct, prod)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) ActivateProduct(prod paywall.Product) error {
	_, err := env.db.NamedExec(paywall.StmtActivateProduct, prod)

	if err != nil {
		return err
	}

	return nil
}

// ListProducts list all products, with each product's plans attached without discount.
// A product's plans are retrieve using JSON_ARRAYAGG.
func (env Env) ListProducts() ([]paywall.ListedProduct, error) {
	products := make([]paywall.ListedProduct, 0)

	err := env.db.Select(&products, paywall.StmtListPricedProducts)

	if err != nil {
		return products, err
	}

	return products, nil
}
