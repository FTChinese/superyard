package admin

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/reader"
)

// FtcAccount retrieves an ftc account before granting/revoking vip.
func (env Env) FtcAccount(ftcID string) (reader.FtcAccount, error) {
	var a reader.FtcAccount
	err := env.db.Get(&a, reader.StmtFtcAccount, ftcID)
	if err != nil {
		return a, err
	}

	return a, nil
}

func (env Env) ListVIP(p gorest.Pagination) ([]reader.FtcAccount, error) {
	var vips = make([]reader.FtcAccount, 0)
	err := env.db.Select(&vips, reader.StmtListVIP, p.Limit, p.Offset())
	if err != nil {
		return nil, err
	}

	return vips, nil
}

// UpdateVIP set/removes vip column.
func (env Env) UpdateVIP(a reader.FtcAccount) error {
	_, err := env.db.NamedExec(reader.StmtSetVIP, a)

	if err != nil {
		return err
	}

	return nil
}
