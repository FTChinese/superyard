package reader

import (
	"errors"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/models/validator"
	"time"
)

// GenerateMemberID generates a random string to membership id.
func GenerateMemberID() string {
	return "mmb_" + rand.String(12)
}

// Membership contains a user's membership information
type Membership struct {
	ID null.String `json:"id" db:"member_id"`
	AccountID
	LegacyTier    null.Int       `json:"-" db:"vip_type"`
	LegacyExpire  null.Int       `json:"-" db:"expire_time"`
	Tier          enum.Tier      `json:"tier" db:"tier"`
	Cycle         enum.Cycle     `json:"cycle" db:"cycle"`
	ExpireDate    chrono.Date    `json:"expireDate" db:"expire_date"`
	PaymentMethod enum.PayMethod `json:"paymentMethod" db:"payment_method"`
	StripeSubID   null.String    `json:"stripeSubId" db:"stripe_sub_id"`
	StripePlanID  null.String    `json:"stripePlanId" db:"stripe_plan_id"`
	AutoRenewal   bool           `json:"autoRenewal" db:"auto_renewal"`
	Status        SubStatus      `json:"status" db:"sub_status"`
}

// NewMember creates a membership directly for a user.
// This is currently used by activating gift cards.
// If membership is purchased via direct payment channel,
// membership is created from subscription order.
func NewMember(accountID AccountID) Membership {
	return Membership{
		ID:        null.StringFrom(GenerateMemberID()),
		AccountID: accountID,
	}
}

func (m *Membership) GenerateID() {
	m.ID = null.StringFrom(GenerateMemberID())
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

func (m Membership) Validate() *validator.InputError {
	if m.Tier == enum.TierNull {
		return &validator.InputError{
			Message: "tier must be one of 'standard' or 'premium'",
			Field:   "tier",
			Code:    validator.CodeInvalid,
		}
	}

	if m.Tier == enum.TierPremium && m.Cycle == enum.CycleMonth {
		return &validator.InputError{
			Message: "monthly subscription is not provided to premium membership",
			Field:   "cycle",
			Code:    validator.CodeInvalid,
		}
	}

	if m.Cycle == enum.CycleNull {
		r := view.NewReason()
		r.SetMessage("cycle must be one of 'month' or 'year'")
		r.Field = "cycle"
		r.Code = view.CodeInvalid

		return &validator.InputError{
			Message: "cycle must be one of 'month' or 'year'",
			Field:   "cycle",
			Code:    validator.CodeInvalid,
		}
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

// FromAliOrWx builds a new membership based on a confirmed
// order.
func (m Membership) FromAliOrWx(sub Order) (Membership, error) {
	if !sub.IsConfirmed() {
		return m, errors.New("only confirmed order could be used to build membership")
	}

	if m.ID.IsZero() {
		m.GenerateID()
	}

	if m.IsZero() {
		m.CompoundID = sub.CompoundID
		m.FtcID = sub.FtcID
		m.UnionID = sub.UnionID
	}

	m.Tier = sub.Tier
	m.Cycle = sub.Cycle
	m.ExpireDate = sub.EndDate
	m.PaymentMethod = sub.PaymentMethod
	m.StripeSubID = null.String{}
	m.StripePlanID = null.String{}
	m.AutoRenewal = false

	return m, nil
}
