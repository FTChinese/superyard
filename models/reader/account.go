package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// FtcID is used to identify an FTC user.
// A user might have an ftc uuid, or a wechat union id,
// or both.
// This type structure is used to ensure unique constraint
// for SQL columns that cannot be both null since SQL do not
// have a mechanism to do UNIQUE INDEX on two columns while
// keeping either of them nullable.
// A user's compound id is taken from either ftc uuid or
// wechat id, with ftc id taking precedence.
type AccountID struct {
	CompoundID string      `json:"-"`
	FtcID      null.String `json:"ftcId"`
	UnionID    null.String `json:"unionId"`
}

// User contains the minimal information to identify a user.
type User struct {
	UserID   string      `json:"id"`
	UnionID  null.String `json:"unionId"`
	Email    string      `json:"email"`
	UserName null.String `json:"userName"`
	IsVIP    bool        `json:"isVip"`
}

// Account show the essential information of a ftc user.
// Client might show a list of accounts and uses those data to query a user's profile, orders, etc.
type Account struct {
	User
	Mobile     null.String `json:"mobile"`
	Nickname   null.String `json:"nickname"`
	Membership Membership  `json:"membership"`
	CreatedAt  chrono.Time `json:"createdAt"`
	UpdatedAt  chrono.Time `json:"updatedAt"`
}
