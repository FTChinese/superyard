package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/models/util"
)

type WxInfo struct {
	UnionID  string      `json:"unionId" db:"union_id"`
	Nickname null.String `json:"nickname" db:"nickname"`
}

// WxUser shows a wechat user's bare-bone data in
// search result.
type WxAccount struct {
	WxInfo
	CreatedAt chrono.Time `json:"createdAt" db:"created_at"`
	UpdatedAt chrono.Time `json:"updatedAt" db:"updated_at"`
	FtcID     null.String `json:"ftcId" db:"ftc_id"`
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
