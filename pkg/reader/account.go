package reader

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/guregu/null"
	"strings"
)

// Wechat contain the essential data to identify a wechat user.
type Wechat struct {
	WxNickname  null.String `json:"nickname" db:"wx_nickname"`
	WxAvatarURL null.String `json:"avatarUrl" db:"wx_avatar_url"`
}

type FtcAccount struct {
	FtcID    null.String `json:"ftcId" db:"ftc_id"`
	UnionID  null.String `json:"unionId" db:"union_id"`
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

// FtcWxAccount contains both ftc cols and wechat cols
// Mainly used as search result.
type FtcWxAccount struct {
	FtcAccount
	Wechat Wechat      `json:"wechat"`
	Kind   AccountKind `json:"kind"`
}

func (a *FtcWxAccount) SetKind() {
	if a.FtcID.Valid {
		a.Kind = AccountKindFtc
		return
	}

	a.Kind = AccountKindWx
}

// Account contains a complete user account, consisting of
// both ftc account and wechat account.
type Account struct {
	FtcWxAccount
	Membership subs.Membership `json:"membership"`
}

type AccountSchema struct {
	FtcAccount
	Wechat
	VIP bool `db:"is_vip"`
	Err error
}

func (s AccountSchema) FtcWxAccount() FtcWxAccount {
	a := FtcWxAccount{
		FtcAccount: s.FtcAccount,
		Wechat:     s.Wechat,
	}

	a.SetKind()

	return a
}

func (s AccountSchema) BuildAccount(m subs.Membership) Account {
	if s.VIP {
		m.Tier = enum.TierVIP
	}

	return Account{
		FtcWxAccount: s.FtcWxAccount(),
		Membership:   m,
	}
}
