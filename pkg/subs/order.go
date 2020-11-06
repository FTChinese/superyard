package subs

import (
	"errors"
	"fmt"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/guregu/null"
	"strconv"
	"strings"
	"time"
)

type Charge struct {
	// The actual amount payable.
	Amount   float64 `json:"amount" db:"amount"`     // Actual price paid.
	Currency string  `json:"currency" db:"currency"` // in which currency.
}

// AliPrice converts Charged price to ailpay format
func (c Charge) AliPrice() string {
	return strconv.FormatFloat(c.Amount, 'f', 2, 32)
}

// AmountInCent converts Charged price to int64 in cent for comparison with wx notification.
func (c Charge) AmountInCent() int64 {
	return int64(c.Amount * 100)
}

func (c Charge) ReadableAmount() string {
	return fmt.Sprintf("%s%.2f",
		strings.ToUpper(c.Currency),
		c.Amount,
	)
}

// Order is a user's subs order
type Order struct {
	ID    string  `json:"id" db:"order_id"`
	Price float64 `json:"price" db:"price"`
	Charge
	CompoundID string      `json:"compoundId" db:"compound_id"`
	FtcID      null.String `json:"ftcId" db:"ftc_id"`
	UnionID    null.String `json:"unionId" db:"union_id"`
	PlanID     null.String `json:"planId" db:"plan_id"`
	DiscountID null.String `json:"discountId" db:"discount_id"`
	paywall.Edition
	Currency      null.String    `json:"currency"`
	CycleCount    int64          `json:"cycleCount" db:"cycle_count"`
	ExtraDays     int64          `json:"extraDays" db:"extra_days"`
	Kind          enum.OrderKind `json:"kind" db:"kind"`
	PaymentMethod enum.PayMethod `json:"payMethod" db:"payment_method"`
	TotalBalance  null.Float     `json:"totalBalance" db:"total_balance"`
	WxAppID       null.String    `json:"wxAppId" db:"wx_app_id"` // Wechat specific. Used by webhook to verify notification.
	CreatedAt     chrono.Time    `json:"createdAt" db:"created_utc"`
	ConfirmedAt   chrono.Time    `json:"confirmedAt" db:"confirmed_utc"`
	StartDate     chrono.Date    `json:"startDate" db:"start_date"`
	EndDate       chrono.Date    `json:"endDate" db:"end_date"`
}

func (o Order) IsConfirmed() bool {
	return !o.ConfirmedAt.IsZero()
}

// Select a time whichever comes later among order's confirmation time and membership's
// expiration time.
func (o Order) pickStartDate(expireDate chrono.Date) chrono.Date {
	if o.Kind == enum.OrderKindUpgrade || o.ConfirmedAt.Time.After(expireDate.Time) {
		return chrono.DateFrom(o.ConfirmedAt.Time)
	}

	return expireDate
}

func (o Order) getEndDate() (chrono.Date, error) {
	var endTime time.Time

	switch o.Cycle {
	case enum.CycleYear:
		endTime = o.StartDate.AddDate(int(o.CycleCount), 0, int(o.ExtraDays))

	case enum.CycleMonth:
		endTime = o.StartDate.AddDate(0, int(o.CycleCount), int(o.ExtraDays))

	default:
		return chrono.Date{}, errors.New("invalid billing cycle")
	}

	return chrono.DateFrom(endTime), nil
}

// Confirm confirms an order based on existing membership.
// If current membership is not expired, the order's
// purchased start date starts from the membership's
// expiration date; otherwise it starts from the
// confirmation time received by webhook.
// If this order is used for upgrading, it always starts
// at now.
func (o Order) Confirm(m reader.Membership) (Order, error) {

	confirmedAt := time.Now()
	o.ConfirmedAt = chrono.TimeFrom(confirmedAt)
	o.StartDate = o.pickStartDate(m.ExpireDate)

	endDate, err := o.getEndDate()
	if err != nil {
		return o, err
	}

	o.EndDate = endDate

	return o, nil
}

// Membership build a membership based on this order.
// The order must be already confirmed.
func (o Order) Membership() (reader.Membership, error) {
	if !o.IsConfirmed() {
		return reader.Membership{}, fmt.Errorf("order %s used to build membership is not confirmed yet", o.ID)
	}

	return reader.Membership{
		CompoundID: null.StringFrom(o.CompoundID),
		IDs: reader.IDs{
			FtcID:   o.FtcID,
			UnionID: o.UnionID,
		},
		LegacyTier:   null.Int{},
		LegacyExpire: null.Int{},
		Edition:      o.Edition,
		ExpireDate:   o.EndDate,
		PayMethod:    o.PaymentMethod,
		FtcPlanID:    o.PlanID,
		StripeSubsID: null.String{},
		StripePlanID: null.String{},
		AutoRenewal:  false,
		Status:       enum.SubsStatusNull,
		AppleSubsID:  null.String{},
		B2BLicenceID: null.String{},
	}, nil
}
