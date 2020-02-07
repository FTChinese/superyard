package search

import (
	"gitlab.com/ftchinese/superyard/models/builder"
	"gitlab.com/ftchinese/superyard/models/reader"
)

func (env Env) SearchFtcUser(email string) (reader.FtcInfo, error) {
	var i reader.FtcInfo

	err := env.DB.Get(&i, stmtSearchFtc, email)

	if err != nil {
		logger.WithField("trace", "Env.SearchFtcUser").Error(err)

		return i, err
	}

	return i, nil
}

func (env Env) SearchWxUser(nickname string, p builder.Pagination) ([]reader.WxInfo, error) {
	wx := make([]reader.WxInfo, 0)

	err := env.DB.Select(
		&wx,
		stmtSearchWx,
		nickname,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.SearchWxUser").Error(err)
		return nil, err
	}

	return wx, nil
}
