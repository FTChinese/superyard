package test

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"time"
)

func NewPaywallBanner() paywall.Banner {
	SeedGoFake()

	return paywall.NewBanner(paywall.BannerInput{
		Heading:    gofakeit.Sentence(10),
		CoverURL:   null.StringFrom(gofakeit.URL()),
		SubHeading: null.StringFrom(gofakeit.Sentence(5)),
		Content:    null.StringFrom(gofakeit.Paragraph(3, 2, 5, "\n")),
	}, gofakeit.Username())
}

func NewPaywallPeriod() paywall.Period {
	return paywall.Period{
		StartUTC: chrono.TimeNow(),
		EndUTC:   chrono.TimeFrom(time.Now().AddDate(0, 0, 1)),
	}
}

func NewPaywallPromo() paywall.Promo {
	SeedGoFake()

	return paywall.NewPromo(paywall.PromoInput{
		BannerInput: paywall.BannerInput{
			Heading:    gofakeit.Sentence(10),
			CoverURL:   null.StringFrom(gofakeit.URL()),
			SubHeading: null.StringFrom(gofakeit.Sentence(5)),
			Content:    null.StringFrom(gofakeit.Paragraph(3, 2, 5, "\n")),
		},
		Period: NewPaywallPeriod(),
	}, gofakeit.Username())
}

type ProductMocker struct {
	id      string
	tier    enum.Tier
	creator string
	price   float64
}

func NewProductMocker(t enum.Tier) ProductMocker {
	SeedGoFake()
	var price float64
	switch t {
	case enum.TierStandard:
		price = 258

	case enum.TierPremium:
		price = 1998
	}

	return ProductMocker{
		id:      paywall.GenProductID(),
		tier:    t,
		creator: gofakeit.Username(),
		price:   price,
	}
}

func (m ProductMocker) Product() paywall.Product {
	SeedGoFake()

	return paywall.Product{
		ID: m.id,
		ProductInput: paywall.ProductInput{
			Tier:        m.tier,
			Heading:     gofakeit.Word(),
			Description: null.StringFrom(gofakeit.Paragraph(4, 1, 5, "\n")),
			SmallPrint:  null.StringFrom(gofakeit.Sentence(10)),
		},
		IsActive:   false,
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
		CreatedBy:  m.creator,
	}
}

func (m ProductMocker) Plan(c enum.Cycle) paywall.Plan {
	SeedGoFake()

	input := paywall.PlanInput{
		ProductID:   m.id,
		Price:       m.price,
		Cycle:       c,
		Description: null.String{},
	}

	if m.tier == enum.TierPremium {
		input.Cycle = enum.CycleYear
		input.Price = 28
	}

	return m.Product().NewPlan(input, m.creator)
}

func (m ProductMocker) PricedProduct() paywall.PricedProduct {
	plans := []paywall.Plan{
		m.Plan(enum.CycleYear),
	}

	if m.tier == enum.TierStandard {
		plans = append(plans, m.Plan(enum.CycleMonth))
	}

	return paywall.PricedProduct{
		Product: m.Product(),
		Plans:   plans,
	}
}

func NewPeriod() paywall.Period {
	return paywall.Period{
		StartUTC: chrono.TimeNow(),
		EndUTC:   chrono.TimeFrom(time.Now().AddDate(0, 0, 1)),
	}
}

func NewDiscount(plan paywall.Plan) paywall.DiscountSchema {
	input := paywall.DiscountInput{
		PriceOff: null.FloatFrom(59),
		Percent:  null.Int{},
		Period:   NewPeriod(),
	}

	return paywall.NewDiscountSchema(input, plan.ID, plan.CreatedBy)
}
