package reader

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/ids"
	"github.com/guregu/null"
)

// Wechat contain the essential data to identify a wechat user.
type Wechat struct {
	WxNickname  null.String `json:"nickname" db:"wx_nickname"`
	WxAvatarURL null.String `json:"avatarUrl" db:"wx_avatar_url"`
}

// BaseAccount contains ftc-only reader account data.
type BaseAccount struct {
	ids.UserIDs
	StripeID null.String `json:"stripeId" db:"stripe_id"`
	Email    null.String `json:"email" db:"email"`
	UserName null.String `json:"userName" db:"user_name"`
	VIP      bool        `json:"vip" db:"is_vip"`
}

// JoinedAccount contains both ftc account and wechat account.
type JoinedAccount struct {
	BaseAccount
	Wechat Wechat `json:"wechat"`
}

// JoinedAccountSchema is used as SQL scan target to retrieve both ftc account and wechat account in a JOIN statement.
type JoinedAccountSchema struct {
	BaseAccount
	Wechat
}

func (s JoinedAccountSchema) JoinedAccount() JoinedAccount {
	a := JoinedAccount{
		BaseAccount: s.BaseAccount,
		Wechat:      s.Wechat,
	}

	return a
}

// Deprecated
type FtcAccountList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []BaseAccount `json:"data"`
	Err  error         `json:"-"`
}
