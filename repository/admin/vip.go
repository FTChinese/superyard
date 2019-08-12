package admin

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/employee"
)

// ListVIP list all vip account on ftchinese.com
func (env Env) ListVIP(p gorest.Pagination) ([]employee.FtcAccount, error) {

	ftcAccounts := make([]employee.FtcAccount, 0)
	err := env.DB.Select(
		ftcAccounts,
		stmtSelectVIP,
		p.Limit,
		p.Offset(),
	)

	if err != nil {
		logger.WithField("trace", "Env.ListVIP").Error(err)

		return nil, err
	}

	return ftcAccounts, nil
}

func (env Env) updateVIP(ftcID string, isVIP bool) error {

	_, err := env.DB.Exec(stmtUpdateVIP, isVIP, ftcID)

	if err != nil {
		logger.WithField("trace", "Env.updateVIP")

		return err
	}

	return nil
}

// GrantVIP set a ftc account as vip
func (env Env) GrantVIP(ftcID string) error {
	return env.updateVIP(ftcID, true)
}

// RevokeVIP removes vip status from a ftc account
func (env Env) RevokeVIP(ftcID string) error {
	return env.updateVIP(ftcID, false)
}
