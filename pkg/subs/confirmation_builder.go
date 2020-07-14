package subs

import (
	"github.com/guregu/null"
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

func (b *ConfirmationBuilder) Build() (ConfirmationResult, error) {

	order, err := b.order.Confirmed(b.mmb)
	if err != nil {
		return ConfirmationResult{}, nil
	}

	m, err := b.mmb.FromAliOrWx(order)
	if err != nil {
		return ConfirmationResult{}, err
	}

	snapshot := b.mmb.Snapshot(b.order.Kind.SnapshotReason())
	snapshot.OrderID = null.StringFrom(order.ID)

	return ConfirmationResult{
		Order:      order,
		Membership: m.Normalize(),
		Snapshot:   snapshot,
	}, nil
}
