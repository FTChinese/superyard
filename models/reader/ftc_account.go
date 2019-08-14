package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/models/util"
)

type FtcInfo struct {
	ID    string `json:"id" db:"ftc_id"`
	Email string `json:"email" db:"email"`
	IsVIP bool   `json:"isVip" db:"is_vip"`
}

type FtcAccount struct {
	FtcInfo
	UnionID   null.String `json:"unionId" db:"union_id"`
	StripeID  null.String `json:"stripeId" db:"stripe_id"`
	UserName  null.String `json:"userName" db:"user_name"`
	Mobile    null.String `json:"mobile" db:"mobile"`
	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt chrono.Time `json:"updatedAt" db:"updated_at"`
}

// LoginHistory identifies how and from where the user login
type LoginHistory struct {
	UserID     string           `json:"userId" db:"user_id"`
	AuthMethod enum.LoginMethod `json:"loginMethod" db:"login_method"`
	util.ClientApp
	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
}
