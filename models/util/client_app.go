package util

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type ClientApp struct {
	ClientType enum.Platform `json:"clientType" db:"client_type"`
	Version    null.String   `json:"clientVersion" db:"client_version"`
	UserIP     null.String   `json:"userIp" db:"user_ip"`
	UserAgent  null.String   `json:"userAgent" db:"user_agent"`
}
