package paywall

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPlan(t *testing.T) {
	seedGoFake()

	input := PlanInput{
		ProductID:   GenProductID(),
		Price:       258,
		Tier:        enum.TierStandard,
		Cycle:       enum.CycleYear,
		Description: null.StringFrom(gofakeit.Paragraph(4, 1, 5, "\n")),
	}

	plan := NewPlan(input, gofakeit.Username())

	assert.NotEmpty(t, plan.ID)

	t.Logf("%s", mustStringify(plan))
}

func TestDiscountedPlanSchema_DiscountedPlan(t *testing.T) {
	schema := DiscountedPlanSchema{
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

	t.Logf("Plan with discount %+v", schema.DiscountedPlan())

	t.Logf("%s", mustStringify(schema.DiscountedPlan()))
}
