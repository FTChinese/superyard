package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// Membership contains a user's membership information
type Membership struct {
	AccountID
	Tier          enum.Tier      `json:"tier"`
	Cycle         enum.Cycle     `json:"cycle"`
	ExpireDate    chrono.Date    `json:"expireDate"`
	PaymentMethod enum.PayMethod `json:"payMethod"`
	StripeSubID   null.String    `json:"stripeSubId"`
	StripePlanID  null.String    `json:"stripePlanId"`
	AutoRenewal   bool           `json:"autoRenewal"`
	Status        string         `json:"status"`
}
