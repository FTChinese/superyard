package user

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
)

// Membership contains a user's membership information
type Membership struct {
	Tier       enum.Tier   `json:"tier"`
	Cycle      enum.Cycle  `json:"cycle"`
	ExpireDate chrono.Date `json:"expireDate"`
}
