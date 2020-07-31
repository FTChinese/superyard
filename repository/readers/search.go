package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/pkg/reader"
)

// SearchFtcAccount tries to find an FTC account by email.
// Returns a single account if found, since the email is
// uniquely constrained.
// I tries to use wildcard search on email: `%<email>%`
// but it's really slow. We'll figure out other way for
// fuzzy match.
func (env Env) SearchFtcAccount(email string, p gorest.Pagination) ([]reader.FtcWxAccount, error) {
	var raws = []reader.AccountSchema{}

	err := env.DB.Select(
		&raws,
		reader.StmtSearchFtcByEmail,
		email,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.SearchFtcAccount").Error(err)
		return nil, err
	}

	var accounts = make([]reader.FtcWxAccount, 0)
	for _, raw := range raws {
		accounts = append(accounts, raw.FtcWxAccount())
	}

	return accounts, nil
}

// SearchWxAccounts tries to find out all wechat user with a LIKE statement.
func (env Env) SearchWxAccounts(nickname string, p gorest.Pagination) ([]reader.FtcWxAccount, error) {
	// NOTE: JOSN marshal result for the empty array is `[]`
	// while for `var rawAccounts []reader.FtcAccount` is `null`.
	var rawAccounts = []reader.AccountSchema{}

	err := env.DB.Select(
		&rawAccounts,
		reader.StmtSearchWxByName,
		"%"+nickname+"%",
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.SearchWxAccounts").Error(err)
		return nil, err
	}

	accounts := make([]reader.FtcWxAccount, 0)
	for _, raw := range rawAccounts {
		accounts = append(accounts, raw.FtcWxAccount())
	}

	return accounts, nil
}
