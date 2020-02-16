package readers

import (
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/models/util"
	"strings"
)

func (env Env) SearchFtcIDByEmail(email string) (string, error) {
	var ftcId string

	if err := env.DB.Get(&ftcId, selectFtcIDByEmail, email); err != nil {
		return "", err
	}

	return ftcId, nil
}

func (env Env) SearchWxIDs(nickname string, p util.Pagination) ([]string, error) {
	var ids []string

	err := env.DB.Select(&ids, selectWxIDs, nickname, p.Limit, p.Offset())
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (env Env) RetrieveWxAccounts(ids []string) ([]reader.BaseAccount, error) {
	var accounts []reader.BaseAccount

	err := env.DB.Select(&accounts, selectWxAccounts, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	for i := range accounts {
		accounts[i].Kind = reader.AccountKindWx
	}

	return accounts, nil
}
