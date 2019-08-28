package reader

import (
	"errors"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"strings"
	"time"
)

func GenerateOrderID() (string, error) {

	id, err := gorest.RandomHex(8)
	if err != nil {
		return "", err
	}

	return "FT" + strings.ToUpper(id), nil
}

// Order is a user's subs order
type Order struct {
	ID string `json:"orderId" db:"order_id"`
	AccountID
	Price            float64        `json:"price" db:"price"`
	Amount           float64        `json:"amount" db:"amount"`
	Tier             enum.Tier      `json:"tier" db:"tier"`
	Cycle            enum.Cycle     `json:"cycle" db:"cycle"`
	Currency         null.String    `json:"currency"`
	CycleCount       int64          `json:"cycleCount" db:"cycle_count"`
	ExtraDays        int64          `json:"extraDays" db:"extra_days"`
	Usage            SubsKind       `json:"usage" db:"usage_type"`
	PaymentMethod    enum.PayMethod `json:"payMethod" db:"payment_method"`
	CreatedAt        chrono.Time    `json:"createdAt" db:"created_at"`
	ConfirmedAt      chrono.Time    `json:"confirmedAt" db:"confirmed_at"`
	StartDate        chrono.Date    `json:"startDate" db:"start_date"`
	EndDate          chrono.Date    `json:"endDate" db:"end_date"`
	UpgradeID        null.String    `json:"-" db:"upgrade_id"`
	MemberSnapshotID null.String    `json:"-" db:"member_snapshot_id"`
}

// GenerateID add the ID field of an incoming order.
func (o *Order) GenerateID() error {
	id, err := GenerateOrderID()
	if err != nil {
		return err
	}

	o.ID = id

	return nil
}

func (o Order) IsConfirmed() bool {
	return !o.ConfirmedAt.IsZero()
}

func (o Order) getStartDate(m Membership, confirmedAt time.Time) time.Time {
	var startTime time.Time

	// If membership is expired, always use the confirmation
	// time as start time.
	if m.IsExpired() {
		startTime = confirmedAt
	} else {
		// If membership is not expired, this order might be
		// used to either renew or upgrade.
		// For renewal, we use current membership's
		// expiration date;
		// For upgrade, we use confirmation time.
		if o.Usage == SubsKindUpgrade {
			startTime = confirmedAt
		} else {
			startTime = m.ExpireDate.Time
		}
	}

	return startTime
}

func (o Order) getEndDate(startTime time.Time) (time.Time, error) {
	var endTime time.Time

	switch o.Cycle {
	case enum.CycleYear:
		endTime = startTime.AddDate(int(o.CycleCount), 0, int(o.ExtraDays))

	case enum.CycleMonth:
		endTime = startTime.AddDate(0, int(o.CycleCount), int(o.ExtraDays))

	default:
		return endTime, errors.New("invalid billing cycle")
	}

	return endTime, nil
}

// Confirm updates an order with existing membership.
// Zero membership is a valid value.
func (o Order) Confirm(m Membership, confirmedAt time.Time) (Order, error) {

	startTime := o.getStartDate(m, confirmedAt)
	endTime, err := o.getEndDate(startTime)
	if err != nil {
		return o, err
	}

	o.ConfirmedAt = chrono.TimeFrom(confirmedAt)
	o.StartDate = chrono.DateFrom(startTime)
	o.EndDate = chrono.DateFrom(endTime)

	return o, nil
}
