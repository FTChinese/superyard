package readers

import (
	"database/sql"
	"gitlab.com/ftchinese/superyard/models/reader"
)

// retrieveFTCAccount retrieves account by ftc id
func (env Env) retrieveFTCAccount(ftcID string) (reader.BaseAccount, error) {
	var a reader.BaseAccount

	if err := env.DB.Get(&a, selectAccountByFtcID, ftcID); err != nil {
		return a, err
	}

	a.SetKind()
	return a, nil
}

func (env Env) retrieveFtcMember(ftcID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, memberByCompoundID, ftcID)

	if err != nil && err != sql.ErrNoRows {
		return reader.Membership{}, err
	}

	m.Normalize()

	return m, nil
}

type accountResult struct {
	success reader.BaseAccount
	err     error
}

type memberResult struct {
	success reader.Membership
	err     error
}

func (env Env) LoadFTCAccount(ftcID string) (reader.Account, error) {
	aChan := make(chan accountResult)
	mChan := make(chan memberResult)

	go func() {
		account, err := env.retrieveFTCAccount(ftcID)
		aChan <- accountResult{
			success: account,
			err:     err,
		}
	}()

	go func() {
		member, err := env.retrieveFtcMember(ftcID)
		mChan <- memberResult{
			success: member,
			err:     err,
		}
	}()

	accountResult, memberResult := <-aChan, <-mChan

	if accountResult.err != nil {
		return reader.Account{}, accountResult.err
	}

	// Ignore ErrNoRows since a reader might not have a membership.
	if memberResult.err != nil && memberResult.err != sql.ErrNoRows {
		return reader.Account{}, memberResult.err
	}

	memberResult.success.VIP = accountResult.success.VIP

	return reader.Account{
		BaseAccount: accountResult.success,
		Membership:  memberResult.success,
	}, nil
}

// retrieveWxAccount retrieve account by wxchat union id.
func (env Env) retrieveWxAccount(unionID string) (reader.BaseAccount, error) {
	var a reader.BaseAccount

	if err := env.DB.Get(&a, selectAccountByWxID, unionID); err != nil {
		return a, err
	}

	a.SetKind()
	return a, nil
}

// retrieveWxMember retrieve membership for wechat.
func (env Env) retrieveWxMember(unionID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, memberByUnionID, unionID)
	if err != nil && err != sql.ErrNoRows {
		return reader.Membership{}, err
	}

	m.Normalize()

	return m, nil
}

func (env Env) LoadWxAccount(unionID string) (reader.Account, error) {
	aChan := make(chan accountResult)
	mChan := make(chan memberResult)

	go func() {
		a, err := env.retrieveWxAccount(unionID)
		aChan <- accountResult{
			success: a,
			err:     err,
		}
	}()

	go func() {
		m, err := env.retrieveWxMember(unionID)
		mChan <- memberResult{
			success: m,
			err:     err,
		}
	}()

	accountResult, memberResult := <-aChan, <-mChan

	if accountResult.err != nil {
		return reader.Account{}, accountResult.err
	}

	if memberResult.err != nil && memberResult.err != sql.ErrNoRows {
		return reader.Account{}, memberResult.err
	}

	memberResult.success.VIP = accountResult.success.VIP

	return reader.Account{
		BaseAccount: accountResult.success,
		Membership:  memberResult.success,
	}, nil
}
