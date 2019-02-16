package user

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

// WxInfo contains a wechat user's information
type WxInfo struct {
	UnionID    string      `json:"unionid"`
	Nickname   string      `json:"nickname"`
	AvatarURL  string      `json:"headimgurl"`
	Gender     enum.Gender `json:"gender"` // 1 for male, 2 for female, 0 for not set.
	Country    string      `json:"country"`
	Province   string      `json:"province"`
	City       string      `json:"city"`
	Privileges []string    `json:"privilege"`
	CreatedAt  chrono.Time `json:"createdAt"`
	UpdatedAt  chrono.Time `json:"updatedAt"`
}

type OAuthHistory struct {
	UnionID string `json:"unionId"`
	OpenID  string `json:"openid"`
	AppID   string `json:"appId"`
	ClientApp
	CreatedAt chrono.Time
	UpdatedAt chrono.Time
}
