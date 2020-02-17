package readers

import (
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/models/util"
)

func (env Env) ListActivities(ftcID string, p util.Pagination) ([]reader.Activity, error) {
	var activities []reader.Activity

	err := env.DB.Select(&activities, stmtActivity, ftcID, p.Limit, p.Offset())

	if err != nil {
		logger.WithField("trace", "ListActivities").Error(err)
		return nil, err
	}

	return activities, nil
}

// ListWxLoginHistory shows a wechat user's login history.
func (env Env) ListWxLoginHistory(unionID string, p util.Pagination) ([]reader.OAuthHistory, error) {

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
