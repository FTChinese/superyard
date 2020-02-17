package readers

import (
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/models/util"
	"strings"
)

func (env Env) findFtcIDByEmail(email string) (string, error) {
	var ftcId string

	if err := env.DB.Get(&ftcId, selectFtcIDByEmail, email); err != nil {
		return "", err
	}

	return ftcId, nil
}

func (env Env) SearchFtcAccount(email string) (reader.BaseAccount, error) {
	ftcID, err := env.findFtcIDByEmail(email)
	if err != nil {
		return reader.BaseAccount{}, err
	}

	a, err := env.retrieveFTCAccount(ftcID)
	if err != nil {
		return reader.BaseAccount{}, err
	}

	return a, nil
}

func (env Env) findWxIDs(nickname string, p util.Pagination) ([]string, error) {
	var ids []string

	err := env.DB.Select(&ids, selectWxIDs, nickname, p.Limit, p.Offset())
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (env Env) retrieveWxAccounts(ids []string) ([]reader.BaseAccount, error) {
	var accounts []reader.BaseAccount

	err := env.DB.Select(&accounts, selectWxAccounts, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	for _, v := range accounts {
		v.SetKind()
	}

	return accounts, nil
}

func (env Env) SearchWxAccounts(nickname string, p util.Pagination) ([]reader.BaseAccount, error) {
	unionIDs, err := env.findWxIDs(nickname, p)
	if err != nil {
		return nil, err
	}

	accounts, err := env.retrieveWxAccounts(unionIDs)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
