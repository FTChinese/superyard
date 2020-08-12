package products

import "github.com/FTChinese/superyard/pkg/paywall"

func (env Env) ListPricedProducts() ([]paywall.PricedProduct, error) {
	schema := make([]paywall.PricedProductSchema, 0)
	products := make([]paywall.PricedProduct, 0)

	err := env.db.Select(&schema, paywall.StmtListPricedProducts)

	if err != nil {
		return products, err
	}

	for _, v := range schema {
		prod, err := v.PricedProduct()
		if err != nil {
			return products, err
		}

		products = append(products, prod)
	}

	return products, nil
}

func (env Env) CreatePricedProduct(p paywall.PricedProduct) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(paywall.StmtCreateProduct, p)

	if err != nil {
		getLogger("CreatePricedProduct").Error(err)
		return err
	}

	if len(p.Plans) > 0 {
		for _, v := range p.Plans {
			_, err := tx.NamedExec(paywall.StmtCreatePlan, v)
			return err
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
