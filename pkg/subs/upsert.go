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

type FtcSubsInput struct {
	reader.IDs
	paywall.Edition
	ExpireDate chrono.Date    `json:"expireDate"`
	PayMethod  enum.PayMethod `json:"payMethod"`
}

func (i *FtcSubsInput) Validate() *render.ValidationError {
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

func (i FtcSubsInput) Membership(a reader.JoinedAccount, plan paywall.Plan) reader.Membership {
	return reader.Membership{
		CompoundID:   null.StringFrom(a.MustGetCompoundID()),
		IDs:          a.IDs,
		LegacyTier:   null.Int{},
		LegacyExpire: null.Int{},
		Edition: paywall.Edition{
			Tier:  plan.Tier,
			Cycle: plan.Cycle,
		},
		ExpireDate:   i.ExpireDate,
		PayMethod:    i.PayMethod,
		FtcPlanID:    null.StringFrom(plan.ID),
		StripeSubsID: null.String{},
		StripePlanID: null.String{},
		AutoRenewal:  false,
		Status:       0,
		AppleSubsID:  null.String{},
		B2BLicenceID: null.String{},
	}
}
