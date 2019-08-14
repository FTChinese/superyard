package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type AccountID struct {
	CompoundID string      `json:"-"`
	FtcID      null.String `json:"ftcId" db:"ftc_id"`
	UnionID    null.String `json:"unionId" db:"union_id"`
}

// Membership contains a user's membership information
type Membership struct {
	AccountID
	Tier          enum.Tier      `json:"tier" db:"tier"`
	Cycle         enum.Cycle     `json:"cycle" db:"cycle"`
	ExpireDate    chrono.Date    `json:"expireDate" db:"expire_date"`
	PaymentMethod enum.PayMethod `json:"paymentMethod" db:"payment_method"`
	StripeSubID   null.String    `json:"stripeSubId" db:"stripe_sub_id"`
	StripePlanID  null.String    `json:"stripePlanId" db:"stripe_plan_id"`
	AutoRenewal   bool           `json:"autoRenewal" db:"auto_renewal"`
	Status        string         `json:"status" db:"sub_status"`
}
