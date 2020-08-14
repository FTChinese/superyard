package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

type GiftCard struct {
	ID         int64       `json:"id"`
	Serial     string      `json:"serial"`
	ExpireDate chrono.Date `json:"expireDate"`
	RedeemedAt chrono.Time `json:"redeemedAt"`
	Tier       enum.Tier   `json:"tier"`
	CycleUnit  enum.Cycle  `json:"cycleUnit"`
	CycleCount int64       `json:"cycleCount"`
}
