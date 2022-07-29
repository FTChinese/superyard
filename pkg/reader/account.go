package reader

import (
	"github.com/FTChinese/superyard/pkg/ids"
	"github.com/guregu/null"
)

// BaseAccount contains ftc-only reader account data.
type BaseAccount struct {
	ids.UserIDs
	StripeID   null.String `json:"stripeId" db:"stripe_id"`
	Email      null.String `json:"email" db:"email"`
	UserName   null.String `json:"userName" db:"user_name"`
	WxNickname null.String `json:"nickname" db:"wx_nickname"`
	VIP        bool        `json:"vip" db:"is_vip"`
}
