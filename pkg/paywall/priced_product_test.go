package paywall

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPricedProduct_Validate(t *testing.T) {
	product := NewPricedProduct(PricedProductInput{
		ProductInput: ProductInput{
			Tier:        enum.TierStandard,
			Heading:     "Standard Edition",
			Description: null.StringFrom("Test"),
			SmallPrint:  null.String{},
		},
		Plans: []PlanInput{
			{
				Cycle:       enum.CycleYear,
				Description: null.StringFrom("Yearly"),
				Price:       0,
			},
		},
	}, "weiguo.ni")

	ve := product.Validate()
	t.Logf("Field: %s", ve.Field)
	assert.Equal(t, ve.Field, "plans.0.price")
}
