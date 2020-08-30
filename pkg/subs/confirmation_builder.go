package subs

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/reader"
)

// ConfirmationBuilder is used to confirm an order, update it based on existing membership expiration date, and then
// update existing membership to next billing cycle.
// This is only used to handle Alipay or Wechat pay.
type ConfirmationBuilder struct {
	mmb   reader.Membership
	order Order
}

func NewConfirmationBuilder(o Order, m reader.Membership) *ConfirmationBuilder {
	return &ConfirmationBuilder{
		mmb:   m,
		order: o,
	}
}

func (b *ConfirmationBuilder) Validate() error {
	if b.order.IsConfirmed() {
		return ErrAlreadyConfirmed
	}

	// If the membership is not expired and created from payment method other than wxpay or alipay,
	// we should deny such requests.
	if !b.mmb.IsExpired() && !b.mmb.IsAliOrWxPay() {
		return ErrValidNonAliOrWxPay
	}

	// If membership is already premium edition.
	if b.order.Kind == enum.OrderKindUpgrade && b.mmb.Tier == enum.TierPremium {
		return ErrAlreadyUpgraded
	}

	return nil
}

func (b *ConfirmationBuilder) Build() (ConfirmationResult, error) {

	order, err := b.order.Confirm(b.mmb)
	if err != nil {
		return ConfirmationResult{}, nil
	}

	m, err := order.Membership()
	if err != nil {
		return ConfirmationResult{}, err
	}

	return ConfirmationResult{
		Order:      order,
		Membership: m.Normalize(),
		Snapshot:   reader.NewSnapshot(reader.SnapshotReasonForOrder(order.Kind), b.mmb).WithOrderID(order.ID),
	}, nil
}
