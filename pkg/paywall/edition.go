package paywall

import "github.com/FTChinese/go-rest/enum"

type Edition struct {
	Tier  enum.Tier  `json:"tier" db:"tier"`
	Cycle enum.Cycle `json:"cycle" db:"cycle"`
}
