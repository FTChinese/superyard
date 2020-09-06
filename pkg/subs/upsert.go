package subs

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strings"
)

type FtcSubsInput struct {
	CompoundID string           `json:"-"`
	Kind       enum.AccountKind `json:"kind"`
	PlanID     string           `json:"planId"`
	PayMethod  enum.PayMethod   `json:"payMethod"`
	ExpireDate chrono.Date      `json:"expireDate"`
}

func (i *FtcSubsInput) Validate() *render.ValidationError {
	i.PlanID = strings.TrimSpace(i.PlanID)

	if i.Kind != enum.AccountKindFtc && i.Kind != enum.AccountKindWx {
		return &render.ValidationError{
			Message: "Account kind is required",
			Field:   "kind",
			Code:    render.CodeMissingField,
		}
	}

	if i.PayMethod != enum.PayMethodAli && i.PayMethod != enum.PayMethodWx {
		return &render.ValidationError{
			Message: "Payment method must be one of alipay or wxpay",
			Field:   "payMethod",
			Code:    render.CodeInvalid,
		}
	}

	return validator.New("planId").Required().Validate(i.PlanID)
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
