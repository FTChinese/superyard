package user

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/types/util"
)

type OAuthAccess struct {
	SessionID string
	// Example: 16_Ix0E3WfWs9u5Rh9f-lB7_LgsQJ4zm1eodolFJpSzoQibTAuhIlp682vDmkZSaYIjD9gekOa1zQl-6c6S_CrN_cN9vx9mybwXNVgFbwPMMwM
	AccessToken string `json:"access_token"`
	// Example: 7200
	ExpiresIn int64 `json:"expires_in"`
	// Exmaple: 16_IlmA9eLGjJw7gBKBT48wff1V1hAYAdpmIqUAypspepm6DsQ6kkcLeZmP932s9PcKp1WM5P_1YwUNQqF-29B_0CqGTqMpWkaaiNSYp26MmB4
	RefreshToken string `json:"refresh_token"`
	// Example: ob7fA0h69OO0sTLyQQpYc55iF_P0
	OpenID string `json:"openid"`
	// Example: snsapi_userinfo
	Scope string `json:"scope"`
	// Example: String:ogfvwjk6bFqv2yQpOrac0J3PqA0o Valid:true
	UnionID null.String `json:"unionid"`
}

// OAuthHistory is a record every time user logged in
// vai Wechat.
type OAuthHistory struct {
	UnionID string `json:"unionId"`
	OpenID  string `json:"openId"`
	AppID   string `json:"appId"`
	util.ClientApp
	CreatedAt chrono.Time `json:"createdAt"`
	UpdatedAt chrono.Time `json:"updatedAt"`
}

// WxInfo contains a wechat user's information.
// This type exists for testing purpose. It is not meant to
// be used by client.
type WxInfo struct {
	UnionID    string      `json:"unionId"`
	Nickname   string      `json:"nickname"`
	AvatarURL  string      `json:"avatarUrl"`
	Gender     enum.Gender `json:"gender"` // 1 for male, 2 for female, 0 for not set.
	Country    string      `json:"country"`
	Province   string      `json:"province"`
	City       string      `json:"city"`
	Privileges []string    `json:"privileges"`
	CreatedAt  chrono.Time `json:"createdAt"`
	UpdatedAt  chrono.Time `json:"updatedAt"`
}
