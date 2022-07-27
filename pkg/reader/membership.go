package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/ids"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/guregu/null"
	"time"
)

var tierToCode = map[enum.Tier]int64{
	enum.TierStandard: 10,
	enum.TierPremium:  100,
}

var codeToTier = map[int64]enum.Tier{
	10:  enum.TierStandard,
	100: enum.TierPremium,
}

// Membership contains a user's membership information
// Creation/Updating strategy:
// Use `PayMethod` to determine how to create/update membership.
// * `Alipay` or `Wechat`: client should manually specify Tier, Cycle, ExpireDate, and StripeSubsID, StripePlanID,
// AutoRenewal, Status, AppleSubsID, B2BLicenceID should not exist;
// `Stripe`: client should only provide the StripeSubsID and we shall ask Stripe API to find out subscription status;
// `Apple`: client should provide Apple subscription id and we shall ask IAP to find out subscription status;
// `B2B`: client should provide the B2B licence id and we shall check DB to find out the subscription status.
// Even for Alipay and Wechat, we still recommend against modifying the data directly. You should find out a buyer's
// order and see if it is confirmed or not. It it is not confirmed yet, confirm that order and the membership
// will be created/updated accordingly.
type Membership struct {
	CompoundID null.String `json:"compoundId" db:"compound_id"`
	ids.UserIDs
	LegacyTier   null.Int `json:"-" db:"vip_type"`
	LegacyExpire null.Int `json:"-" db:"expire_time"`
	paywall.Edition
	ExpireDate   chrono.Date     `json:"expireDate" db:"expire_date"`
	PayMethod    enum.PayMethod  `json:"payMethod" db:"pay_method"`
	FtcPlanID    null.String     `json:"ftcPlanId" db:"ftc_plan_id"`
	StripeSubsID null.String     `json:"stripeSubsId" db:"stripe_subs_id"` // If it exists, client should refresh.
	StripePlanID null.String     `json:"stripePlanId" db:"stripe_plan_id"`
	AutoRenewal  bool            `json:"autoRenewal" db:"auto_renewal"`
	Status       enum.SubsStatus `json:"status" db:"subs_status"`
	AppleSubsID  null.String     `json:"appleSubsId" db:"apple_subs_id"`   // If exists, client should refresh
	B2BLicenceID null.String     `json:"b2bLicenceId" db:"b2b_licence_id"` // If exists, client should refresh
}

func (m Membership) isLegacyOnly() bool {
	if m.LegacyExpire.Valid && m.LegacyTier.Valid && m.ExpireDate.IsZero() && m.Tier == enum.TierNull {
		return true
	}

	return false
}

func (m Membership) isAPIOnly() bool {
	if (m.LegacyExpire.IsZero() && m.LegacyTier.IsZero()) && (!m.ExpireDate.IsZero() && m.Tier != enum.TierNull) {
		return true
	}

	return false
}

// IsZero test whether the instance is empty.
func (m Membership) IsZero() bool {
	return m.CompoundID.IsZero()
}

// Normalize turns legacy vip_type and expire_time into
// member_tier and expire_date columns, or vice versus.
// Issues: if we set expiration date to an earlier time, data become inconsistent.
func (m Membership) Normalize() Membership {
	if m.IsZero() {
		return m
	}

	// Syn from legacy format to api created columns
	if m.isLegacyOnly() {
		// Note the conversion is not exactly the same moment since Golang converts Unix in local time.
		expireDate := time.Unix(m.LegacyExpire.Int64, 0)

		m.ExpireDate = chrono.DateFrom(expireDate)
		m.Tier = codeToTier[m.LegacyTier.Int64]
		// m.Cycle cannot be determined

		return m
	}

	// Sync from api columns to legacy column
	if m.isAPIOnly() {
		m.LegacyExpire = null.IntFrom(m.ExpireDate.Unix())
		m.LegacyTier = null.IntFrom(tierToCode[m.Tier])

		return m
	}

	// Otherwise do not touch it.
	return m
}
