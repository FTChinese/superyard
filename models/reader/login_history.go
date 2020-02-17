package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/util"
)

type ActivityKind string

const (
	ActivityKindNull          ActivityKind = ""
	ActivityKindLogin                      = "login"
	ActivityKindSignUp                     = "signup"
	ActivityKindVerification               = "email_verification"
	ActivityKindPasswordReset              = "password_reset"
)

// Activity shows a user's footprint when using email account.
type Activity struct {
	FtcID      string        `json:"ftcId" db:"ftc_id"`
	Platform   enum.Platform `json:"platform" db:"platform"`
	Version    null.String   `json:"version" db:"version"`
	UserIP     null.String   `json:"userIp" db:"user_ip"`
	UserAgent  null.String   `json:"userAgent" db:"user_agent"`
	CreatedUTC chrono.Time   `json:"createdUtc" db:"created_utc"`
	Kind       ActivityKind  `json:"kind" db:"kind"`
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
