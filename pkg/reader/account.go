package reader

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"strings"
)

// Wechat contain the essential data to identify a wechat user.
type Wechat struct {
	WxNickname  null.String `json:"nickname" db:"wx_nickname"`
	WxAvatarURL null.String `json:"avatarUrl" db:"wx_avatar_url"`
}

type FtcAccount struct {
	IDs
	StripeID null.String `json:"stripeId" db:"stripe_id"`
	Email    null.String `json:"email" db:"email"`
	UserName null.String `json:"userName" db:"user_name"`
}

func (a FtcAccount) NormalizedName() string {
	if a.UserName.Valid {
		return a.UserName.String
	}

	if a.Email.Valid {
		return strings.Split(a.Email.String, "@")[0]
	}

	return ""
}

// JoinedAccount contains both ftc cols and wechat cols
type JoinedAccount struct {
	FtcAccount
	Wechat Wechat           `json:"wechat"`
	Kind   enum.AccountKind `json:"kind"`
}

func (a *JoinedAccount) SetKind() {
	if a.FtcID.Valid {
		a.Kind = enum.AccountKindFtc
		return
	}

	a.Kind = enum.AccountKindWx
}

// Account contains a complete user account, consisting of
// both ftc account and wechat account.
type Account struct {
	JoinedAccount
	Membership Membership `json:"membership"`
}

type AccountSchema struct {
	FtcAccount
	Wechat
	VIP bool `db:"is_vip"`
	Err error
}

func (s AccountSchema) FtcWxAccount() JoinedAccount {
	a := JoinedAccount{
		FtcAccount: s.FtcAccount,
		Wechat:     s.Wechat,
	}

	a.SetKind()

	return a
}

func (s AccountSchema) BuildAccount(m Membership) Account {
	if s.VIP {
		m.Tier = enum.TierVIP
	}

	return Account{
		JoinedAccount: s.FtcWxAccount(),
		Membership:    m,
	}
}
