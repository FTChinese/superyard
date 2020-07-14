package subs

import (
	"errors"
	"github.com/FTChinese/go-rest/chrono"
	"time"
)

// ConfirmationBuilder is used to confirm an order, update it based on existing membership expiration date, and then
// update existing membership to next billing cycle.
// This is only used to handle Alipay or Wechat pay.
type ConfirmationBuilder struct {
	mmb   Membership
	order Order
}

func NewConfirmationBuilder(o Order, m Membership) *ConfirmationBuilder {
	return &ConfirmationBuilder{
		mmb:   m,
		order: o,
	}
}

// pickStartTime selects the appropriate starting time when
// confirming an order
func (b *ConfirmationBuilder) pickStartTime() time.Time {
	// If current membership is expires, or not exist.
	if b.mmb.IsExpired() {
		return time.Now()
	}

	// For upgrading, it always starts at confirmation time.
	if b.order.Kind == KindUpgrade {
		return time.Now()
	}

	// If membership is not expired, use its expiration date
	return b.mmb.ExpireDate.Time
}

func (b *ConfirmationBuilder) confirmedOrder() (Order, error) {
	startTime := b.pickStartTime()

	endTime, err := b.order.getEndDate(startTime)
	if err != nil {
		return Order{}, errors.New("cannot determine order's end time")
	}

	order := b.order

	order.ConfirmedAt = chrono.TimeNow()
	order.StartDate = chrono.DateFrom(startTime)
	order.EndDate = chrono.DateFrom(endTime)

	return order, nil
}
