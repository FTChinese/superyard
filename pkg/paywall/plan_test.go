package paywall

import (
	"github.com/guregu/null"
	"testing"
)

func TestDiscountedPlanSchema_DiscountedPlan(t *testing.T) {
	schema := ExpandedPlanSchema{
		Plan: planStdYear.Plan,
		Discount: Discount{
			DiscID: null.String{},
			DiscountInput: DiscountInput{
				PriceOff: null.Float{},
				Percent:  null.Int{},
				Period:   Period{},
			},
		},
	}

	t.Logf("Plan with discount %+v", schema.ExpandedPlan())

	t.Logf("%s", mustStringify(schema.ExpandedPlan()))
}
