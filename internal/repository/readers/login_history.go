package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) ListActivities(ftcID string, p gorest.Pagination) ([]reader.Activity, error) {
	var activities []reader.Activity

	err := env.dbs.Read.Select(&activities, reader.StmtActivity, ftcID, p.Limit, p.Offset())

	if err != nil {
		return nil, err
	}

	return activities, nil
}

// ListWxLoginHistory shows a wechat user's login history.
func (env Env) ListWxLoginHistory(unionID string, p gorest.Pagination) ([]reader.OAuthHistory, error) {

	var ah []reader.OAuthHistory

	err := env.dbs.Read.Select(
		&ah,
		reader.StmtWxLoginHistory,
		unionID,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return ah, nil
}
