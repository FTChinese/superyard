package subs

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/guregu/null"
	"strings"
)

type FtcSubsUpdateInput struct {
	paywall.Edition
	ExpireDate chrono.Date    `json:"expireDate"`
	PayMethod  enum.PayMethod `json:"payMethod"`
	PlanID     string         `json:"-"` // Not part of the request body.
}

func (i FtcSubsUpdateInput) Validate() *render.ValidationError {

	if i.Tier == enum.TierNull {
		return &render.ValidationError{
			Message: "Tier is required",
			Field:   "tier",
			Code:    render.CodeMissingField,
		}
	}

	if i.Cycle == enum.CycleNull {
		return &render.ValidationError{
			Message: "Cycle is required",
			Field:   "cycle",
			Code:    render.CodeMissingField,
		}
	}

	if i.Tier == enum.TierPremium && i.Cycle == enum.CycleMonth {
		return &render.ValidationError{
			Message: "Premium edition does not have monthly billing cycle",
			Field:   "cycle",
			Code:    render.CodeInvalid,
		}
	}

	if i.PayMethod != enum.PayMethodAli && i.PayMethod != enum.PayMethodWx {
		return &render.ValidationError{
			Message: "Payment method must be one of alipay or wxpay",
			Field:   "payMethod",
			Code:    render.CodeInvalid,
		}
	}

	return nil
}

type FtcSubsCreationInput struct {
	reader.IDs
	FtcSubsUpdateInput
}

func (i *FtcSubsCreationInput) Validate() *render.ValidationError {
	ftcID := strings.TrimSpace(i.FtcID.String)
	unionID := strings.TrimSpace(i.UnionID.String)

	if ftcID == "" && unionID == "" {
		return &render.ValidationError{
			Message: "Provide at least one of ftc id or wechat union id.",
			Field:   "compoundId",
			Code:    render.CodeMissingField,
		}
	}

	i.FtcID = null.NewString(ftcID, ftcID != "")
	i.UnionID = null.NewString(unionID, unionID != "")

	return i.FtcSubsUpdateInput.Validate()
}

func ManualCreateMember(a reader.JoinedAccount, i FtcSubsCreationInput) reader.Membership {
	return reader.Membership{
		CompoundID:   null.StringFrom(a.MustGetCompoundID()),
		IDs:          a.IDs,
		LegacyTier:   null.Int{},
		LegacyExpire: null.Int{},
		Edition: paywall.Edition{
			Tier:  i.Tier,
			Cycle: i.Cycle,
		},
		ExpireDate:   i.ExpireDate,
		PayMethod:    i.PayMethod,
		FtcPlanID:    null.StringFrom(i.PlanID),
		StripeSubsID: null.String{},
		StripePlanID: null.String{},
		AutoRenewal:  false,
		Status:       0,
		AppleSubsID:  null.String{},
		B2BLicenceID: null.String{},
	}
}

func ManualUpdateMember(m reader.Membership, input FtcSubsUpdateInput) reader.Membership {
	return reader.Membership{
		CompoundID:   m.CompoundID,
		IDs:          m.IDs,
		LegacyTier:   null.Int{},
		LegacyExpire: null.Int{},
		Edition:      input.Edition,
		ExpireDate:   input.ExpireDate,
		PayMethod:    input.PayMethod,
		FtcPlanID:    null.StringFrom(input.PlanID),
		StripeSubsID: null.String{},
		StripePlanID: null.String{},
		AutoRenewal:  false,
		Status:       0,
		AppleSubsID:  null.String{},
		B2BLicenceID: null.String{},
	}
}
