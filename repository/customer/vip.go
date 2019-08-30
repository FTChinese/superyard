package customer

import (
	"github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

// ListVIP list all vip account on ftchinese.com
func (env Env) ListVIP(p gorest.Pagination) ([]reader.FtcInfo, error) {

	// Ignore Goland warning here. We want to send back empty array
	// to indicate no element exists, rather than `null`
	var info = []reader.FtcInfo{}

	err := env.DB.Select(
		&info,
		stmtSelectVIP,
		p.Limit,
		p.Offset(),
	)

	if err != nil {
		logger.WithField("trace", "Env.ListVIP").Error(err)

		return nil, err
	}

	return info, nil
}

// GrantVIP set a ftc account as vip
func (env Env) GrantVIP(ftcID string) error {

	_, err := env.DB.Exec(stmtGrantVIP, ftcID)
	if err != nil {
		return err
	}

	return nil
}

// RevokeVIP removes vip status from a ftc account
func (env Env) RevokeVIP(ftcID string) error {

	_, err := env.DB.Exec(stmtRevokeVIP, ftcID)
	if err != nil {
		return err
	}

	return nil
}
