package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"gitlab.com/ftchinese/superyard/models/util"
)

// LoginHistory identifies how and from where the user login
type LoginHistory struct {
	UserID     string           `json:"userId" db:"user_id"`
	AuthMethod enum.LoginMethod `json:"loginMethod" db:"login_method"`
	util.ClientApp
	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
}

// OAuthHistory is a record every time user logged in
// vai WxAccount.
type OAuthHistory struct {
	UnionID string `json:"unionId" db:"union_id"`
	OpenID  string `json:"openId" db:"open_id"`
	AppID   string `json:"appId" db:"app_id"`
	util.ClientApp
	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt chrono.Time `json:"updatedAt" db:"updated_at"`
}
