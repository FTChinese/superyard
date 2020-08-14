package paywall

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPricedProduct(t *testing.T) {
	seedGoFake()

	input := PricedProductInput{
		ProductInput: ProductInput{
			Tier:        enum.TierStandard,
			Heading:     "Standard Edition",
			Description: null.StringFrom(gofakeit.Paragraph(4, 1, 5, "\n")),
			SmallPrint:  null.String{},
		},
		Plans: []PlanInput{
			{
				Price:       258,
				Cycle:       enum.CycleYear,
				Description: null.StringFrom(gofakeit.Word()),
			},
			{
				Price:       28,
				Cycle:       enum.CycleMonth,
				Description: null.StringFrom(gofakeit.Word()),
			},
		},
	}

	pp := NewPricedProduct(input, gofakeit.Username())

	assert.Len(t, pp.Plans, 2)

	t.Logf("Priced product: %s", mustStringify(input))
}
