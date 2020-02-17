package reader

import (
	"github.com/guregu/null"
)

type BaseAccount struct {
	FtcID    null.String `json:"ftcId" db:"ftc_id"`
	UnionID  null.String `json:"unionId" db:"union_id"`
	StripeID null.String `json:"stripeId" db:"stripe_id"`
	Email    null.String `json:"email" db:"email"`
	UserName null.String `json:"userName" db:"user_name"`
	Nickname null.String `json:"nickname" db:"nickname"`
	Kind     AccountKind `json:"kind"`
}

func (a *BaseAccount) SetKind() {
	if a.FtcID.Valid {
		a.Kind = AccountKindFtc
		return
	}

	a.Kind = AccountKindWx
}

// Account contains a complete user account, consisting of
// both ftc account and wechat account.
type Account struct {
	BaseAccount
	Membership Membership `json:"membership"`
}
