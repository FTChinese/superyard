package test

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/guregu/null"
	"time"
)

var PlanStdYear = paywall.ExpandedPlan{
	Plan: paywall.Plan{
		ID: "plan_MynUQDQY1TSQ",
		PlanInput: paywall.PlanInput{
			Cycle:       enum.CycleYear,
			Description: null.String{},
			Price:       258,
			ProductID:   "prod_zjWdiTUpDN8l",
		},
		Tier:       enum.TierStandard,
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: paywall.Discount{
		DiscID:     null.StringFrom("dsc_F7gEwjaF3OsR"),
		DiscPlanID: null.StringFrom("plan_MynUQDQY1TSQ"),
		DiscountInput: paywall.DiscountInput{
			PriceOff: null.FloatFrom(130),
			Percent:  null.Int{},
			Period: paywall.Period{
				StartUTC: chrono.TimeNow(),
				EndUTC:   chrono.TimeFrom(time.Now().AddDate(0, 0, 2)),
			},
		},
	},
}

var PlanStdMonth = paywall.ExpandedPlan{
	Plan: paywall.Plan{
		ID: "plan_1Uz4hrLy3Mzy",
		PlanInput: paywall.PlanInput{
			Cycle:       enum.CycleMonth,
			Description: null.String{},
			Price:       28,
			ProductID:   "prod_zjWdiTUpDN8l",
		},
		Tier:       enum.TierStandard,
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: paywall.Discount{},
}

var PlanPrm = paywall.ExpandedPlan{
	Plan: paywall.Plan{
		ID: "plan_vRUzRQ3aglea",
		PlanInput: paywall.PlanInput{
			Cycle:       enum.CycleYear,
			Description: null.String{},
			Price:       1998,
			ProductID:   "prod_IaoK5SbK79g8",
		},
		Tier:       enum.TierPremium,
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: paywall.Discount{},
}
