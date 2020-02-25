package readers

import (
	"strings"

	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/models/util"
)

func (env Env) findFtcIDByEmail(email string) (string, error) {
	var ftcID string

	if err := env.DB.Get(&ftcID, selectFtcIDByEmail, email); err != nil {
		logger.WithField("trace", "Env.findFtcIDByEmail").Error(err)
		return "", err
	}

	return ftcID, nil
}

// SearchFtcAccount tries to find an FTC account by email.
// Returns a single account if found, since the email is
// uniquely constrained.
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

// findWxIDs collects all wechat union ids whose nickname contains the
// parameter `nickname`.
func (env Env) findWxIDs(nickname string, p util.Pagination) ([]string, error) {
	var ids []string

	err := env.DB.Select(
		&ids,
		selectWxIDs,
		nickname,
		p.Limit,
		p.Offset())
	// NOTE: ErrNoRows won't be thrown if no rows found for statement selecting
	// multiple rows.
	if err != nil {
		logger.WithField("trace", "Env.findWxIDs").Error(err)
		return nil, err
	}

	// An emtpy slice will be returned if nothing found.
	return ids, nil
}

// retrieveWxAccounts loads wechat acocunt for the passed in union ids.
func (env Env) retrieveWxAccounts(ids []string) ([]reader.BaseAccount, error) {
	// NOTE: JOSN marshal result for the empty array is `[]`
	// while for `var accounts []reader.BaseAccount` is `null`.
	var accounts = []reader.BaseAccount{}

	if len(ids) == 0 {
		return accounts, nil
	}

	err := env.DB.Select(&accounts, selectWxAccounts, strings.Join(ids, ","))
	if err != nil {
		logger.WithField("trace", "Env.retrieveWxAccounts").Error(err)
		return nil, err
	}

	for i := range accounts {
		accounts[i].SetKind()
	}

	return accounts, nil
}

// SearchWxAccounts tries to find out all wechat user with a LIKE statement.
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
