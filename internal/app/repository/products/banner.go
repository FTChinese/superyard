package products

import "github.com/FTChinese/superyard/pkg/paywall"

func (env Env) CreateBanner(b paywall.Banner) error {
	_, err := env.dbs.Write.NamedExec(paywall.StmtCreateBanner, b)

	if err != nil {
		return err
	}

	return nil
}

// LoadBanner retrieves a Banner by id. The id is always 1.
func (env Env) LoadBanner(id int64) (paywall.Banner, error) {
	var b paywall.Banner
	err := env.dbs.Read.Get(&b, paywall.StmtBanner, id)
	if err != nil {
		return b, err
	}

	return b, err
}

func (env Env) UpdateBanner(b paywall.Banner) error {
	_, err := env.dbs.Write.NamedExec(paywall.StmtUpdateBanner, b)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) DropBannerPromo(bannerID int64) error {
	_, err := env.dbs.Write.Exec(paywall.StmtDropPromo, bannerID)
	if err != nil {
		return err
	}

	return nil
}
