package customer

import (
	"database/sql"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

// RetrieveAccountFtc retrieves account by ftc id
func (env Env) RetrieveAccountFtc(ftcID string) (reader.Account, error) {
	var a reader.Account

	if err := env.DB.Get(&a, stmtFtcJoinWx, ftcID); err != nil {
		return reader.Account{}, err
	}

	return a, nil
}

// RetrieveAccountWx retrieve account by wxchat union id.
func (env Env) RetrieveAccountWx(unionID string) (reader.Account, error) {
	var a reader.Account

	if err := env.DB.Get(&a, stmtWxJoinFtc, unionID); err != nil {
		return reader.Account{}, err
	}

	return a, nil
}

func (env Env) RetrieveMemberFtc(ftcID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, memberForEmail, ftcID)

	if err != nil && err != sql.ErrNoRows {
		return reader.Membership{}, err
	}

	m.Normalize()

	return m, nil
}

// RetrieveMemberWx retrieve membership for wechat.
func (env Env) RetrieveMemberWx(unionID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, memberForWx, unionID)
	if err != nil && err != sql.ErrNoRows {
		return reader.Membership{}, err
	}

	m.Normalize()

	return m, nil
}

// RetrieveFtcProfile loads profile of an email user.
func (env Env) RetrieveFtcProfile(ftcID string) (reader.FtcProfile, error) {
	var p reader.FtcProfile

	if err := env.DB.Get(&p, selectFtcProfile, ftcID); err != nil {
		return reader.FtcProfile{}, err
	}

	return p, nil
}

// RetrieveWxProfile loads profile of a wx user.
func (env Env) RetrieveWxProfile(unionID string) (reader.WxProfile, error) {
	var p reader.WxProfile

	if err := env.DB.Get(&p, selectWxProfile, unionID); err != nil {
		return reader.WxProfile{}, err
	}

	return p, nil
}
