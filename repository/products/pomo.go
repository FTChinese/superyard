package products

import (
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/guregu/null"
)

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
