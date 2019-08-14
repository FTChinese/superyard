package customer

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

func (env Env) RetrieveWxAccount(unionID string) (reader.WxAccount, error) {
	var w reader.WxAccount

	err := env.DB.Get(&w, stmtWxAccount, unionID)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveWxAccount").Error(err)

		return w, err
	}

	return w, nil
}

// ListOAuthHistory shows a wechat user's login history.
func (env Env) ListOAuthHistory(unionID string, p gorest.Pagination) ([]reader.OAuthHistory, error) {

	var ah []reader.OAuthHistory

	err := env.DB.Select(
		&ah,
		stmtWxLoginHistory,
		unionID,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "ListOAuthHistory").Error(err)
		return nil, err
	}

	return ah, nil
}
