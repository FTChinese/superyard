package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/reader"
)

// SearchJoinedAccountEmail tries to find an FTC account by email.
// Returns a single account if found, since the email is
// uniquely constrained.
// I tries to use wildcard search on email: `%<email>%`
// but it's really slow. We'll figure out other way for
// fuzzy match.
func (env Env) SearchJoinedAccountEmail(email string, p gorest.Pagination) ([]reader.JoinedAccount, error) {
	var raws = make([]reader.JoinedAccountSchema, 0)

	err := env.dbs.Read.Select(
		&raws,
		reader.StmtSearchJoinedAccountByEmail,
		email,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	var accounts = make([]reader.JoinedAccount, 0)
	for _, raw := range raws {
		accounts = append(accounts, raw.JoinedAccount())
	}

	return accounts, nil
}

// SearchJoinedAccountWxName tries to find out all wechat user with a LIKE statement.
func (env Env) SearchJoinedAccountWxName(nickname string, p gorest.Pagination) ([]reader.JoinedAccount, error) {
	// NOTE: JOSN marshal result for the empty array is `[]`
	// while for `var rawAccounts []reader.FtcAccount` is `null`.
	var rawAccounts = make([]reader.JoinedAccountSchema, 0)

	err := env.dbs.Read.Select(
		&rawAccounts,
		reader.StmtSearchJoinedAccountByWxName,
		"%"+nickname+"%",
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	accounts := make([]reader.JoinedAccount, 0)
	for _, raw := range rawAccounts {
		accounts = append(accounts, raw.JoinedAccount())
	}

	return accounts, nil
}
