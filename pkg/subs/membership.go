package subs

import (
	"errors"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"time"
)

// Membership contains a user's membership information
// Creation/Updating strategy:
// Use `PaymentMethod` to determine how to create/update membership.
// * `Alipay` or `Wechat`: client should manually specify Tier, Cycle, ExpireDate, and StripeSubsID, StripePlanID,
// AutoRenewal, Status, AppleSubsID, B2BLicenceID should not exist;
// `Stripe`: client should only provide the StripeSubsID and we shall ask Stripe API to find out subscription status;
// `Apple`: client should provide Apple subscription id and we shall ask IAP to find out subscription status;
// `B2B`: client should provide the B2B licence id and we shall check DB to find out the subscription status.
// Even for Alipay and Wechat, we still recommend against modifying the data directly. You should find out a buyer's
// order and see if it is confirmed or not. It it is not confirmed yet, confirm that order and the membership
// will be created/updated accordingly.
// TODO: add FTC plan id.
type Membership struct {
	CompoundID    string          `json:"compoundId" db:"compound_id"`
	FtcID         null.String     `json:"ftcId" db:"ftc_id"`
	UnionID       null.String     `json:"unionId" db:"union_id"`
	LegacyTier    null.Int        `json:"-" db:"vip_type"`
	LegacyExpire  null.Int        `json:"-" db:"expire_time"`
	Tier          enum.Tier       `json:"tier" db:"tier"`
	Cycle         enum.Cycle      `json:"cycle" db:"cycle"`
	ExpireDate    chrono.Date     `json:"expireDate" db:"expire_date"`
	PaymentMethod enum.PayMethod  `json:"payMethod" db:"payment_method"`
	StripeSubsID  null.String     `json:"stripeSubsId" db:"stripe_subs_id"` // If it exists, client should refresh.
	StripePlanID  null.String     `json:"stripePlanId" db:"stripe_plan_id"`
	AutoRenewal   bool            `json:"autoRenewal" db:"auto_renewal"`
	Status        enum.SubsStatus `json:"status" db:"subs_status"`
	AppleSubsID   null.String     `json:"appleSubsId" db:"apple_subs_id"`   // If exists, client should refresh
	B2BLicenceID  null.String     `json:"b2bLicenceId" db:"b2b_licence_id"` // If exists, client should refresh
}

// Normalize turns legacy vip_type and expire_time into
// member_tier and expire_date columns, or vice versus.
func (m *Membership) Normalize() {
	// Turn unix seconds to time.
	if m.LegacyExpire.Valid && m.ExpireDate.IsZero() {
		m.ExpireDate = chrono.DateFrom(time.Unix(m.LegacyExpire.Int64, 0))
	}

	// Turn time to unix seconds.
	if !m.ExpireDate.IsZero() && m.LegacyExpire.IsZero() {
		m.LegacyExpire = null.IntFrom(m.ExpireDate.Unix())
	}

	if m.LegacyTier.Valid && m.Tier == enum.TierNull {
		switch m.LegacyTier.Int64 {
		case 10:
			m.Tier = enum.TierStandard
		case 100:
			m.Tier = enum.TierPremium
		}
	}

	if m.Tier != enum.TierNull && m.LegacyTier.IsZero() {
		switch m.Tier {
		case enum.TierStandard:
			m.LegacyTier = null.IntFrom(10)
		case enum.TierPremium:
			m.LegacyTier = null.IntFrom(100)
		}
	}
}

func (m Membership) Validate() *render.ValidationError {
	if m.Tier == enum.TierNull {
		return &render.ValidationError{
			Message: "tier must be one of 'standard' or 'premium'",
			Field:   "tier",
			Code:    render.CodeInvalid,
		}
	}

	if m.Tier == enum.TierPremium && m.Cycle == enum.CycleMonth {
		return &render.ValidationError{
			Message: "monthly subscription is not provided to premium membership",
			Field:   "cycle",
			Code:    render.CodeInvalid,
		}
	}

	if m.Cycle == enum.CycleNull {
		r := view.NewReason()
		r.SetMessage("cycle must be one of 'month' or 'year'")
		r.Field = "cycle"
		r.Code = view.CodeInvalid

		return &render.ValidationError{
			Message: "cycle must be one of 'month' or 'year'",
			Field:   "cycle",
			Code:    render.CodeInvalid,
		}
	}

	if m.PaymentMethod == enum.PayMethodNull {
		return &render.ValidationError{
			Message: "You must specify a payment method",
			Field:   "payMethod",
			Code:    render.CodeMissingField,
		}
	}

	// TODO: ensure fields mutual exclusive.
	if m.PaymentMethod == enum.PayMethodAli || m.PaymentMethod == enum.PayMethodWx {

	}

	return nil
}

// IsZero test whether the instance is empty.
func (m Membership) IsZero() bool {
	return m.CompoundID == "" && m.Tier == enum.TierNull
}

// IsExpired tests if the membership's expiration date is before now.
func (m Membership) IsExpired() bool {
	// If membership does not exist, it is treated as expired.
	if m.IsZero() {
		return true
	}

	// If expire date is before now, AND auto renew is false,
	// we treat this one as actually expired.
	// If ExpireDate is passed, but auto renew is true, we still
	// treat this one as not expired.
	return m.ExpireDate.Before(time.Now().Truncate(24*time.Hour)) && !m.AutoRenewal
}

// FromAliOrWx builds a new membership based on a confirmed order.
// This is used when we are confirming an order.
func (m Membership) FromAliOrWx(order Order) (Membership, error) {
	if !order.IsConfirmed() {
		return m, errors.New("only confirmed order could be used to build membership")
	}

	if m.IsZero() {
		m.CompoundID = order.CompoundID
		m.FtcID = order.FtcID
		m.UnionID = order.UnionID
	}

	m.Tier = order.Tier
	m.Cycle = order.Cycle
	m.ExpireDate = order.EndDate
	m.PaymentMethod = order.PaymentMethod
	m.StripeSubsID = null.String{}
	m.StripePlanID = null.String{}
	m.AutoRenewal = false
	m.AppleSubsID = null.String{}
	m.B2BLicenceID = null.String{}

	return m, nil
}

func (m Membership) Snapshot(reason enum.SnapshotReason) MemberSnapshot {
	return MemberSnapshot{
		ID:         "snp_" + rand.String(12),
		Reason:     reason,
		CreatedUTC: chrono.TimeNow(),
		Membership: m,
	}
}
