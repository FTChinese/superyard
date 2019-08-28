package customer

import (
	"database/sql"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

type ftcAccountResult struct {
	success reader.FtcAccount
	err     error
}

type memberResult struct {
	success reader.Membership
	err     error
}

type wxAccountResult struct {
	success reader.WxAccount
	err     error
}

// LoadAccountFtc retrieves ftc account.
// This is an experimental demo of using multiple
// simple queries to mimic SQL JOIN backed by
// goroutine concurrency.
func (env Env) LoadAccountFtc(ftcID string) (reader.Account, error) {

	ac := make(chan ftcAccountResult)
	mc := make(chan memberResult)

	go func() {
		a, err := env.RetrieveFtcAccount(ftcID)

		ac <- ftcAccountResult{
			success: a,
			err:     err,
		}
	}()
	go func() {
		m, err := env.RetrieveMember(ftcID)

		mc <- memberResult{
			success: m,
			err:     err,
		}
	}()

	accountResult, memberResult := <-ac, <-mc

	if accountResult.err != nil {
		return reader.Account{}, accountResult.err
	}

	if memberResult.err != nil && memberResult.err != sql.ErrNoRows {
		return reader.Account{}, memberResult.err
	}

	ftcAccount := accountResult.success

	account := reader.Account{
		Ftc:        ftcAccount,
		Membership: memberResult.success,
	}

	if ftcAccount.UnionID.Valid {
		wechat, err := env.RetrieveWxAccount(ftcAccount.UnionID.String)
		if err != nil && err != sql.ErrNoRows {
			return reader.Account{}, err
		}

		account.Wechat = wechat
	}

	return account, nil
}

func (env Env) LoadAccountWx(unionID string) (reader.Account, error) {

	wc := make(chan wxAccountResult)
	mc := make(chan memberResult)

	go func() {
		w, err := env.RetrieveWxAccount(unionID)

		wc <- wxAccountResult{
			success: w,
			err:     err,
		}
	}()
	go func() {
		m, err := env.RetrieveMember(unionID)

		mc <- memberResult{
			success: m,
			err:     err,
		}
	}()

	wxResult, memberResult := <-wc, <-mc

	if wxResult.err != nil {
		return reader.Account{}, wxResult.err
	}

	if memberResult.err != nil && memberResult.err != sql.ErrNoRows {
		return reader.Account{}, memberResult.err
	}

	account := reader.Account{
		Ftc:        reader.FtcAccount{},
		Wechat:     wxResult.success,
		Membership: memberResult.success,
	}

	ftcID := wxResult.success.FtcID
	if ftcID.Valid {
		ftcAccount, err := env.RetrieveFtcAccount(ftcID.String)
		if err != nil && err != sql.ErrNoRows {
			return reader.Account{}, err
		}

		account.Ftc = ftcAccount
	}

	return account, nil
}
