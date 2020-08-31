package readers

import (
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) FtcAccount(id string) (reader.FtcAccount, error) {
	var a reader.FtcAccount

	if err := env.db.Get(&a, reader.StmtFtcBaseAccount, id); err != nil {
		return a, err
	}

	return a, nil
}

// accountByFtcID retrieves account by ftc id
func (env Env) accountByFtcID(ftcID string) (reader.JoinedAccountSchema, error) {
	var a reader.JoinedAccountSchema

	if err := env.db.Get(&a, reader.StmtAccountByFtcID, ftcID); err != nil {
		return a, err
	}

	return a, nil
}

// accountByWxID retrieve account by wxchat union id.
func (env Env) accountByWxID(unionID string) (reader.JoinedAccountSchema, error) {
	var a reader.JoinedAccountSchema

	if err := env.db.Get(&a, reader.StmtAccountByWxID, unionID); err != nil {
		return a, err
	}

	return a, nil
}

type accountAsyncResult struct {
	success reader.JoinedAccountSchema
	err     error
}

func (env Env) asyncAccountByFtcID(ftcID string) <-chan accountAsyncResult {
	c := make(chan accountAsyncResult)

	go func() {
		defer close(c)
		a, err := env.accountByFtcID(ftcID)

		c <- accountAsyncResult{
			success: a,
			err:     err,
		}
	}()

	return c
}

func (env Env) asyncAccountByWxID(unionID string) <-chan accountAsyncResult {
	c := make(chan accountAsyncResult)

	go func() {
		defer close(c)
		a, err := env.accountByWxID(unionID)

		c <- accountAsyncResult{
			success: a,
			err:     err,
		}
	}()

	return c
}

func (env Env) AccountByFtcID(ftcID string) (reader.Account, error) {
	aChan, mChan := env.asyncAccountByFtcID(ftcID), env.asyncMembership(ftcID)

	accountResult, memberResult := <-aChan, <-mChan

	if accountResult.err != nil {
		return reader.Account{}, accountResult.err
	}

	// Ignore ErrNoRows since a reader might not have a membership.
	if memberResult.err != nil {
		return reader.Account{}, memberResult.err
	}

	return accountResult.success.BuildAccount(memberResult.value), nil
}

func (env Env) AccountByUnionID(unionID string) (reader.Account, error) {
	aChan, mChan := env.asyncAccountByWxID(unionID), env.asyncMembership(unionID)

	accountResult, memberResult := <-aChan, <-mChan

	if accountResult.err != nil {
		return reader.Account{}, accountResult.err
	}

	if memberResult.err != nil {
		return reader.Account{}, memberResult.err
	}

	return accountResult.success.BuildAccount(memberResult.value), nil
}
