package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/go-rest/view"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strings"
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

// MemberInput specifies the data to create or update membership.
// Using this approach to modify membership data should be avoided
// as possible as you can.
// Membership should only be updated after consulting payment provider.
// For wechant and alipay, use confirm order endpoint;
// For stripe and apple, fetch latest subscription status from them.
type MemberInput struct {
	CompoundID string         `json:"compoundId" db:"compound_id"` // When creating a membership directly, you should provide this value. Use ftc id if present, then fallback to union id.
	ExpireDate chrono.Date    `json:"expireDate" db:"expire_date"`
	PayMethod  enum.PayMethod `json:"payMethod" db:"pay_method"`
	FtcPlanID  null.String    `json:"ftcPlanId" db:"ftc_plan_id"` // Whe use plan id to determine which pricing plan user is subscribed to.
}

func (i MemberInput) NewMembership(a FtcAccount, plan paywall.Plan) Membership {
	return Membership{
		CompoundID:   null.StringFrom(a.MustGetCompoundID()),
		IDs:          a.IDs,
		LegacyTier:   null.Int{},
		LegacyExpire: null.Int{},
		Edition: Edition{
			Tier:  plan.Tier,
			Cycle: plan.Cycle,
		},
		ExpireDate:   i.ExpireDate,
		PayMethod:    i.PayMethod,
		FtcPlanID:    i.FtcPlanID,
		StripeSubsID: null.String{},
		StripePlanID: null.String{},
		AutoRenewal:  false,
		Status:       enum.SubsStatusNull,
		AppleSubsID:  null.String{},
		B2BLicenceID: null.String{},
	}
}

func (i *MemberInput) Validate() *render.ValidationError {
	i.CompoundID = strings.TrimSpace(i.CompoundID)
	i.FtcPlanID.String = strings.TrimSpace(i.FtcPlanID.String)

	if i.PayMethod == enum.PayMethodNull {
		return &render.ValidationError{
			Message: "Payment method is required",
			Field:   "payMethod",
			Code:    render.CodeMissingField,
		}
	}

	if i.PayMethod != enum.PayMethodAli && i.PayMethod != enum.PayMethodWx {
		return &render.ValidationError{
			Message: "It is not supported to manually modify membership with payment method other than alipay or wechat",
			Field:   "payMethod",
			Code:    render.CodeInvalid,
		}
	}

	ve := validator.New("compoundId").Required().Validate(i.CompoundID)
	if ve != nil {
		return ve
	}

	ve = validator.New("ftcPlanId").Required().Validate(i.FtcPlanID.String)
	if ve != nil {
		return ve
	}

	return nil
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
	IDs
	LegacyTier   null.Int `json:"-" db:"vip_type"`
	LegacyExpire null.Int `json:"-" db:"expire_time"`
	Edition
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

// Normalize turns legacy vip_type and expire_time into
// member_tier and expire_date columns, or vice versus.
func (m Membership) Normalize() Membership {
	if m.IsZero() {
		return m
	}

	legacyDate := time.Unix(m.LegacyExpire.Int64, 0)

	// Use whichever comes later.
	// If LegacyExpire is after ExpireDate, then we should
	// use LegacyExpire and LegacyTier
	if legacyDate.After(m.ExpireDate.Time) {
		m.ExpireDate = chrono.DateFrom(legacyDate)
		m.Tier = codeToTier[m.LegacyTier.Int64]
	} else {
		m.LegacyExpire = null.IntFrom(m.ExpireDate.Unix())
		m.LegacyTier = null.IntFrom(tierToCode[m.Tier])
	}

	return m
}

func (m Membership) Update(input MemberInput, plan paywall.Plan) Membership {
	m.ExpireDate = input.ExpireDate
	m.PayMethod = input.PayMethod
	m.FtcPlanID = input.FtcPlanID
	m.Tier = plan.Tier
	m.Cycle = plan.Cycle

	return m
}

// Validate makes sure fields are valid.
// How a membership is created/updated depends on the payment method:
// If payment method == alipay or wecaht, then StripeSubsID, AppleSubsID and B2BLicenceID must
// not exist and the membership is created/updated directly;
// If payment method == stripe and stripe subscription id is provided,
// then fetch this user's subscription data from Stripe and update
// our db accordingly. The data returned from Stripe API is the only source of truth;
// If payment method == apple, then fetch subscription data from IAP, which is the only source of truth;
// If payment method == b2b, then check the b2b licence id status.
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

	if m.PayMethod == enum.PayMethodNull {
		return &render.ValidationError{
			Message: "You must specify a payment method",
			Field:   "payMethod",
			Code:    render.CodeMissingField,
		}
	}

	// TODO: ensure fields mutual exclusive.
	if m.PayMethod != enum.PayMethodAli && m.PayMethod != enum.PayMethodWx {
		return &render.ValidationError{
			Message: "Manually modify membership with payment method other than alipay or wechat is not supported",
			Field:   "payMethod",
			Code:    render.CodeInvalid,
		}
	}

	if m.StripeSubsID.Valid || m.StripePlanID.Valid || m.AutoRenewal || m.Status != enum.SubsStatusNull || m.AppleSubsID.Valid || m.B2BLicenceID.Valid {
		return &render.ValidationError{
			Message: "Manually modify membership with payment method other than alipay or wechat is not supported",
			Field:   "payMethod",
			Code:    render.CodeMissing,
		}
	}

	return nil
}

func (m Membership) ValidateCreate() *render.ValidationError {
	ve := validator.New("compoundId").Required().Validate(m.CompoundID.String)
	if ve != nil {
		return ve
	}

	// FtcID and UnionID cannot be both empty.
	if m.FtcID.IsZero() && m.UnionID.IsZero() {
		ve := validator.New("ftcId").Required().Validate(m.FtcID.String)
		if ve != nil {
			return ve
		}

		ve = validator.New("unionId").Required().Validate(m.UnionID.String)
		if ve != nil {
			return ve
		}
	}

	return m.Validate()
}

// IsZero test whether the instance is empty.
func (m Membership) IsZero() bool {
	return m.CompoundID.IsZero() && m.Tier == enum.TierNull
}

func (m Membership) IsEqual(other Membership) bool {
	return m.CompoundID == other.CompoundID && m.Tier == other.Tier && m.Cycle == other.Cycle && m.PayMethod == other.PayMethod
}

// IsAliOrWxPay checks whether the current membership comes from Alipay or Wxpay.
func (m Membership) IsAliOrWxPay() bool {
	// For backward compatibility. If Tier field comes from LegacyTier, then PayMethod field will be null.
	// We treat all those cases as wxpay or alipay.
	if m.Tier != enum.TierNull && m.PayMethod == enum.PayMethodNull {
		return true
	}

	return m.PayMethod == enum.PayMethodAli || m.PayMethod == enum.PayMethodWx
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
