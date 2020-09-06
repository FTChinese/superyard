package readers

import (
	"database/sql"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/reader"
)

// FtcAccount retrieves ftc-only account by uuid.
func (env Env) FtcAccount(ftcID string) (reader.FtcAccount, error) {
	var a reader.FtcAccount

	if err := env.db.Get(&a, reader.StmtFtcAccount, ftcID); err != nil {
		return a, err
	}

	return a, nil
}

// joinedAccountByFtcID retrieves ftc + wx account by ftc id.
// Wx part might be zero values.
func (env Env) joinedAccountByFtcID(ftcID string) (reader.JoinedAccountSchema, error) {
	var a reader.JoinedAccountSchema

	if err := env.db.Get(
		&a,
		reader.StmtJoinedAccountByFtcID,
		ftcID); err != nil {
		return a, err
	}

	return a, nil
}

// joinedAccountByWxID retrieve ftc + wx account by wxchat union id.
// The ftc part might be zero values.
func (env Env) joinedAccountByWxID(unionID string) (reader.JoinedAccountSchema, error) {
	var a reader.JoinedAccountSchema

	if err := env.db.Get(&a, reader.StmtJoinedAccountByWxID, unionID); err != nil {
		return a, err
	}

	return a, nil
}

func (env Env) JoinedAccountByFtcOrWx(compoundID string, kind enum.AccountKind) (reader.JoinedAccount, error) {
	var schema reader.JoinedAccountSchema
	var err error

	switch kind {
	case enum.AccountKindFtc:
		schema, err = env.joinedAccountByFtcID(compoundID)

	case enum.AccountKindWx:
		schema, err = env.joinedAccountByWxID(compoundID)

	default:
		return reader.JoinedAccount{}, sql.ErrNoRows
	}

	if err != nil {
		return reader.JoinedAccount{}, err
	}

	return schema.JoinedAccount(), err
}

type accountAsyncResult struct {
	success reader.JoinedAccountSchema
	err     error
}

func (env Env) asyncJoinedAccountByFtcID(ftcID string) <-chan accountAsyncResult {
	c := make(chan accountAsyncResult)

	go func() {
		defer close(c)
		a, err := env.joinedAccountByFtcID(ftcID)

		c <- accountAsyncResult{
			success: a,
			err:     err,
		}
	}()

	return c
}

func (env Env) asyncJoinedAccountByWxID(unionID string) <-chan accountAsyncResult {
	c := make(chan accountAsyncResult)

	go func() {
		defer close(c)
		a, err := env.joinedAccountByWxID(unionID)

		c <- accountAsyncResult{
			success: a,
			err:     err,
		}
	}()

	return c
}

func (env Env) AccountByFtcID(ftcID string) (reader.Account, error) {
	aChan, mChan := env.asyncJoinedAccountByFtcID(ftcID), env.asyncMembership(ftcID)

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
	aChan, mChan := env.asyncJoinedAccountByWxID(unionID), env.asyncMembership(unionID)

	accountResult, memberResult := <-aChan, <-mChan

	if accountResult.err != nil {
		return reader.Account{}, accountResult.err
	}

	if memberResult.err != nil {
		return reader.Account{}, memberResult.err
	}

	return accountResult.success.BuildAccount(memberResult.value), nil
}
