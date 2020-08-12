package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var prodStd = Product{
	ID: "prod_oj4ks8shj38",
	ProductInput: ProductInput{
		Tier:        enum.TierStandard,
		Heading:     "标准会员",
		Description: null.StringFrom("专享订阅内容每日仅需{{dailyAverageOfYear}}元(或按月订阅每日{{dailyAverageOfMonth}}元)\r\n精选深度分析\r\n中英双语内容\r\n金融英语速读训练\r\n英语原声电台\r\n无限浏览7日前所有历史文章（近8万篇）"),
		SmallPrint:  null.String{},
	},
	CreatedUTC: chrono.TimeNow(),
	UpdatedUTC: chrono.TimeNow(),
	CreatedBy:  "weiguo.ni",
}

var prodPrm = Product{
	ID: "prod_35rbbrgz08c",
	ProductInput: ProductInput{
		Tier:        enum.TierPremium,
		Heading:     "高端会员",
		Description: null.StringFrom("专享订阅内容每日仅需{{dailyAverageOfYear}}元\r\n享受“标准会员”所有权益\r\n编辑精选，总编/各版块主编每周五为您推荐本周必读资讯，分享他们的思考与观点\r\nFT中文网2018年度论坛门票2张，价值3999元/张 （不含差旅与食宿）"),
		SmallPrint:  null.StringFrom("注：所有活动门票不可折算现金、不能转让、不含差旅与食宿"),
	},
	CreatedUTC: chrono.TimeNow(),
	UpdatedUTC: chrono.TimeNow(),
	CreatedBy:  "weiguo.ni",
}

var planStdYear = DiscountedPlan{
	Plan: Plan{
		ID: "plan_ICMPPM0UXcpZ",
		PlanInput: PlanInput{
			ProductID:   prodStd.ID,
			Price:       258,
			Tier:        enum.TierStandard,
			Cycle:       enum.CycleYear,
			Description: null.StringFrom("Standard monthly price"),
		},
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: Discount{
		ID:     null.StringFrom("dsc_ykDueW15nIJQ"),
		PlanID: null.StringFrom("plan_ICMPPM0UXcpZ"),
		DiscountInput: DiscountInput{
			PriceOff: null.IntFrom(59),
			Percent:  null.Int{},
			Period: Period{
				StartUTC: chrono.TimeNow(),
				EndUTC:   chrono.TimeNow(),
			},
		},
	},
}

var planStdMonth = DiscountedPlan{
	Plan: Plan{
		ID: "plan_wl5esy783d",
		PlanInput: PlanInput{
			ProductID:   prodStd.ID,
			Price:       28,
			Tier:        enum.TierStandard,
			Cycle:       enum.CycleMonth,
			Description: null.StringFrom("Standard monthly price"),
		},
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: Discount{
		ID:     null.String{},
		PlanID: null.String{},
		DiscountInput: DiscountInput{
			PriceOff: null.Int{},
			Percent:  null.Int{},
			Period:   Period{},
		},
	},
}

var planPrmYear = DiscountedPlan{
	Plan: Plan{
		ID: "plan_5iIonqaehig4",
		PlanInput: PlanInput{
			ProductID:   prodPrm.ID,
			Price:       1998,
			Tier:        enum.TierPremium,
			Cycle:       enum.CycleYear,
			Description: null.StringFrom("Premium yearly price"),
		},
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: Discount{
		ID:     null.String{},
		PlanID: null.String{},
		DiscountInput: DiscountInput{
			PriceOff: null.Int{},
			Percent:  null.Int{},
			Period: Period{
				StartUTC: chrono.Time{},
				EndUTC:   chrono.Time{},
			},
		},
	},
}

func TestGroupPlans(t *testing.T) {
	result := GroupPlans([]DiscountedPlan{
		planStdYear,
		planPrmYear,
		planStdMonth,
	})

	assert.Equal(t, len(result), 2)
	assert.Equal(t, len(result[prodStd.ID]), 2)
	assert.Equal(t, len(result[prodPrm.ID]), 1)
}

func TestBuildPaywallProducts(t *testing.T) {
	result := BuildPaywallProducts(
		[]Product{prodStd, prodPrm},
		[]DiscountedPlan{planStdYear, planPrmYear, planStdMonth},
	)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, len(result[0].Plans), 2)
	assert.Equal(t, len(result[1].Plans), 1)
}

func TestNewPromo(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())
	input := PromoInput{
		BannerInput: BannerInput{
			Heading:    gofakeit.Sentence(10),
			CoverURL:   null.StringFrom(gofakeit.URL()),
			SubHeading: null.StringFrom(gofakeit.Sentence(5)),
			Content:    null.StringFrom(gofakeit.Paragraph(3, 2, 5, "\n")),
		},
		Period: Period{
			StartUTC: chrono.TimeNow(),
			EndUTC:   chrono.TimeNow(),
		},
	}

	p := NewPromo(input, "weiguo.ni")

	assert.NotEmpty(t, p.ID)

	t.Logf("%+v", p)
}
