package products

import "github.com/FTChinese/superyard/pkg/paywall"

func (env Env) CreateDiscount(d paywall.DiscountSchema) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(paywall.StmtCreateDiscount, d)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(paywall.StmtApplyDiscount, d)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (env Env) DropDiscount(plan paywall.Plan) error {
	_, err := env.db.NamedExec(paywall.StmtDropDiscount, plan)
	if err != nil {
		return err
	}

	return nil
}
