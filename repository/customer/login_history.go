package customer

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

func (env Env) ListEmailLoginHistory(ftcID string, p gorest.Pagination) ([]reader.LoginHistory, error) {

	var lh []reader.LoginHistory

	err := env.DB.Select(
		&lh,
		stmtLoginHistory,
		ftcID,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "ListEmailLoginHistory").Error(err)
		return nil, err
	}

	return lh, nil
}

// ListWxLoginHistory shows a wechat user's login history.
func (env Env) ListWxLoginHistory(unionID string, p gorest.Pagination) ([]reader.OAuthHistory, error) {

	var ah []reader.OAuthHistory

	err := env.DB.Select(
		&ah,
		stmtWxLoginHistory,
		unionID,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "ListWxLoginHistory").Error(err)
		return nil, err
	}

	return ah, nil
}