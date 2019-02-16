package user

import "github.com/FTChinese/go-rest/enum"

type ClientApp struct {
	ClientType enum.Platform `json:"clientType"`
	Version    null.String `json:"clientVersion"`
	UserIP     null.String `json:"userIp"`
	UserAgent  null.String `json:"userAgent"`
}
