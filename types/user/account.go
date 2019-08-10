package user

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// Membership contains a user's membership information
type Membership struct {
	Tier       enum.Tier   `json:"tier"`
	Cycle      enum.Cycle  `json:"cycle"`
	ExpireDate chrono.Date `json:"expireDate"`
}

// User contains the minimal information to identify a user.
type User struct {
	UserID   string      `json:"id"`
	UnionID  null.String `json:"unionId"`
	Email    string      `json:"email"`
	UserName null.String `json:"userName"`
	IsVIP    bool        `json:"isVip"`
}

// WxUser shows a wechat user's bare-bone data in
// search result.
type Wechat struct {
	UnionID   string      `json:"unionId"`
	Nickname  null.String `json:"nickname"`
	CreatedAt chrono.Time `json:"createdAt"`
	UpdatedAt chrono.Time `json:"updatedAt"`
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