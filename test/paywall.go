package test

import (
	"bytes"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/dt"
	"github.com/FTChinese/superyard/pkg/ids"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/pkg/price"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"io"
	"time"
)

func NewPaywallBanner() paywall.Banner {
	faker.SeedGoFake()

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
	faker.SeedGoFake()

	return paywall.NewPromo(paywall.PromoInput{
		Heading:    null.StringFrom(gofakeit.Sentence(10)),
		CoverURL:   null.StringFrom(gofakeit.URL()),
		SubHeading: null.StringFrom(gofakeit.Sentence(5)),
		Content:    null.StringFrom(gofakeit.Paragraph(3, 2, 5, "\n")),
		Terms:      null.StringFrom(gofakeit.Paragraph(3, 2, 5, "\n")),
		Period:     NewPaywallPeriod(),
	}, gofakeit.Username())
}

type ProductMocker struct {
	id      string
	tier    enum.Tier
	creator string
	price   float64
}

func NewProductMocker(t enum.Tier) ProductMocker {
	faker.SeedGoFake()
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
	faker.SeedGoFake()

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
	faker.SeedGoFake()

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

type ProductBuilder struct {
	productID string
	tier      enum.Tier
}

func NewProductBuilder(id string) ProductBuilder {
	if id == "" {
		id = ids.ProductID()
	}

	return ProductBuilder{
		productID: ids.ProductID(),
		tier:      enum.TierStandard,
	}
}

func (b ProductBuilder) WithStd() ProductBuilder {
	b.tier = enum.TierStandard
	return b
}

func (b ProductBuilder) WithPrm() ProductBuilder {
	b.tier = enum.TierPremium
	return b
}

func (b ProductBuilder) NewPriceBuilder(id string) PriceBuilder {
	if id == "" {
		id = ids.PriceID()
	}
	return PriceBuilder{
		productID: b.productID,
		edition: price.Edition{
			Tier:  b.tier,
			Cycle: enum.CycleYear,
		},
		live: true,
	}
}

type PriceBuilder struct {
	productID string
	edition   price.Edition
	live      bool
}

func (b PriceBuilder) WithYear() PriceBuilder {
	b.edition.Cycle = enum.CycleYear
	return b
}

func (b PriceBuilder) WithMonth() PriceBuilder {
	b.edition.Cycle = enum.CycleMonth
	return b
}

func (b PriceBuilder) WithLive() PriceBuilder {
	b.live = true
	return b
}

func (b PriceBuilder) WithTest() PriceBuilder {
	b.live = false
	return b
}

func (b PriceBuilder) Build() price.FtcPriceParams {
	var amount float64
	if b.edition == price.StdMonthEdition {
		amount = 35
	} else if b.edition == price.StdYearEdition {
		amount = 298
	} else if b.edition == price.PremiumEdition {
		amount = 1998
	}

	return price.FtcPriceParams{
		CreatedBy:   gofakeit.Username(),
		Edition:     b.edition,
		Description: null.String{},
		LiveMode:    b.live,
		Nickname:    null.String{},
		Price:       amount,
		ProductID:   b.productID,
	}
}

func (b PriceBuilder) BuildIOBody() io.Reader {
	body := faker.MustMarshalIndent(b.Build())
	return bytes.NewReader(body)
}

func (b PriceBuilder) NewDiscountBuilder(priceID string) DiscountBuilder {
	var off float64
	if b.edition == price.StdMonthEdition {
		off = 34
	} else if b.edition == price.StdYearEdition {
		off = 50
	} else if b.edition == price.PremiumEdition {
		off = 100
	}

	return DiscountBuilder{
		priceID: priceID,
		off:     off,
	}
}

type DiscountBuilder struct {
	priceID string
	off     float64
}

func (b DiscountBuilder) Build(k price.OfferKind) price.DiscountParams {
	return price.DiscountParams{
		CreatedBy:   gofakeit.Username(),
		Description: null.StringFrom(gofakeit.Sentence(10)),
		Kind:        k,
		Percent:     null.Int{},
		DateTimePeriod: dt.DateTimePeriod{
			StartUTC: chrono.TimeNow(),
			EndUTC:   chrono.TimeFrom(time.Now().AddDate(0, 0, 7)),
		},
		PriceOff:  null.FloatFrom(b.off),
		PriceID:   b.priceID,
		Recurring: false,
	}
}

func (b DiscountBuilder) BuildIntro() price.DiscountParams {
	return b.Build(price.OfferKindIntroductory)
}

func (b DiscountBuilder) BuildPromo() price.DiscountParams {
	return b.Build(price.OfferKindPromotion)
}
func (b DiscountBuilder) BuildRetention() price.DiscountParams {
	return b.Build(price.OfferKindRetention)
}
func (b DiscountBuilder) BuildWinBack() price.DiscountParams {
	return b.Build(price.OfferKindWinBack)
}
