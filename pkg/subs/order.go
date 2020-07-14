package subs

import (
	"errors"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"time"
)

// OrderPeriod is  a duration an order purchased.
type OrderPeriod struct {
	StartDate chrono.Date `json:"startDate" db:"start_date"`
	EndDate   chrono.Date `json:"endDate" db:"end_date"`
}

// Order is a user's subs order
type Order struct {
	ID               string         `json:"id" db:"order_id"`
	CompoundID       string         `json:"compoundId" db:"compound_id"`
	FtcID            null.String    `json:"ftcId" db:"ftc_id"`
	UnionID          null.String    `json:"unionId" db:"union_id"`
	Price            float64        `json:"price" db:"price"`
	Amount           float64        `json:"amount" db:"amount"`
	Tier             enum.Tier      `json:"tier" db:"tier"`
	Cycle            enum.Cycle     `json:"cycle" db:"cycle"`
	Currency         null.String    `json:"currency"`
	CycleCount       int64          `json:"cycleCount" db:"cycle_count"`
	ExtraDays        int64          `json:"extraDays" db:"extra_days"`
	Kind             Kind           `json:"kind" db:"usage_type"`
	PaymentMethod    enum.PayMethod `json:"paymentMethod" db:"payment_method"`
	CreatedAt        chrono.Time    `json:"createdAt" db:"created_at"`
	ConfirmedAt      chrono.Time    `json:"confirmedAt" db:"confirmed_at"`
	StartDate        chrono.Date    `json:"startDate" db:"start_date"`
	EndDate          chrono.Date    `json:"endDate" db:"end_date"`
	UpgradeID        null.String    `json:"-" db:"upgrade_id"`
	MemberSnapshotID null.String    `json:"-" db:"member_snapshot_id"`
}

func (o Order) IsConfirmed() bool {
	return !o.ConfirmedAt.IsZero()
}

func (o Order) startTimeAfterConfirmed(m Membership, confirmedAt time.Time) time.Time {
	// If current membership is expires, or not exist.
	if m.IsExpired() {
		return confirmedAt
	}

	// For upgrading, it always starts at confirmation time.
	if o.Kind == KindUpgrade {
		return confirmedAt
	}

	// If membership is not expired, use its expiration date
	return m.ExpireDate.Time
}

func (o Order) endTimeAfterConfirmed(start time.Time) (time.Time, error) {
	switch o.Cycle {
	case enum.CycleYear:
		return start.AddDate(int(o.CycleCount), 0, int(o.ExtraDays)), nil

	case enum.CycleMonth:
		return start.AddDate(0, int(o.CycleCount), int(o.ExtraDays)), nil

	default:
		return time.Time{}, errors.New("invalid billing cycle")
	}
}

// Confirmed confirms an order based on current membership
// expiration status.
// Zero membership is a valid value.
func (o Order) Confirmed(m Membership) (Order, error) {

	confirmedAt := time.Now()

	startTime := o.startTimeAfterConfirmed(m, confirmedAt)
	endTime, err := o.endTimeAfterConfirmed(startTime)
	if err != nil {
		return o, err
	}

	o.ConfirmedAt = chrono.TimeFrom(confirmedAt)
	o.StartDate = chrono.DateFrom(startTime)
	o.EndDate = chrono.DateFrom(endTime)

	return o, nil
}
