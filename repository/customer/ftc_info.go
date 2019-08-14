package customer

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

func (env Env) RetrieveFtcAccount(ftcID string) (reader.FtcAccount, error) {
	var a reader.FtcAccount

	err := env.DB.Get(&a, stmtFtcInfo, ftcID)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveFtcAccount").Error(err)
		return a, err
	}

	return a, nil
}

func (env Env) ListLoginHistory(ftcID string, p gorest.Pagination) ([]reader.LoginHistory, error) {

	var lh []reader.LoginHistory

	err := env.DB.Select(
		&lh,
		stmtLoginHistory,
		ftcID,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "ListLoginHistory").Error(err)
		return nil, err
	}

	return lh, nil
}
