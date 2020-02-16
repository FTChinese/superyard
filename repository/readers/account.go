package readers

import (
	"database/sql"
	"gitlab.com/ftchinese/superyard/models/reader"
)

// RetrieveAccountFtc retrieves account by ftc id
func (env Env) RetrieveAccountFtc(ftcID string) (reader.BaseAccount, error) {
	var a reader.BaseAccount

	if err := env.DB.Get(&a, selectAccountByFtcID, ftcID); err != nil {
		return a, err
	}

	a.Kind = reader.AccountKindFtc
	return a, nil
}

// RetrieveAccountWx retrieve account by wxchat union id.
func (env Env) RetrieveAccountWx(unionID string) (reader.BaseAccount, error) {
	var a reader.BaseAccount

	if err := env.DB.Get(&a, selectAccountByWxID, unionID); err != nil {
		return a, err
	}

	a.Kind = reader.AccountKindWx
	return a, nil
}

func (env Env) RetrieveMemberFtc(ftcID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, memberByCompoundID, ftcID)

	if err != nil && err != sql.ErrNoRows {
		return reader.Membership{}, err
	}

	m.Normalize()

	return m, nil
}

// RetrieveMemberWx retrieve membership for wechat.
func (env Env) RetrieveMemberWx(unionID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, memberByUnionID, unionID)
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
