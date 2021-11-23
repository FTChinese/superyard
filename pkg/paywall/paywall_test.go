package paywall

import (
	"encoding/json"
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

var planStdYear = ExpandedPlan{
	Plan: Plan{
		ID: "plan_ICMPPM0UXcpZ",
		PlanInput: PlanInput{
			ProductID:   prodStd.ID,
			Price:       258,
			Cycle:       enum.CycleYear,
			Description: null.StringFrom("Standard monthly price"),
		},
		Tier:       enum.TierStandard,
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: Discount{
		DiscID: null.StringFrom("dsc_ykDueW15nIJQ"),
		DiscountInput: DiscountInput{
			PriceOff: null.FloatFrom(59),
			Percent:  null.Int{},
			Period: Period{
				StartUTC: chrono.TimeNow(),
				EndUTC:   chrono.TimeNow(),
			},
		},
	},
}

var planStdMonth = ExpandedPlan{
	Plan: Plan{
		ID: "plan_wl5esy783d",
		PlanInput: PlanInput{
			ProductID:   prodStd.ID,
			Price:       28,
			Cycle:       enum.CycleMonth,
			Description: null.StringFrom("Standard monthly price"),
		},
		Tier:       enum.TierStandard,
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: Discount{
		DiscID: null.String{},
		DiscountInput: DiscountInput{
			PriceOff: null.Float{},
			Percent:  null.Int{},
			Period:   Period{},
		},
	},
}

var planPrmYear = ExpandedPlan{
	Plan: Plan{
		ID: "plan_5iIonqaehig4",
		PlanInput: PlanInput{
			ProductID:   prodPrm.ID,
			Price:       1998,
			Cycle:       enum.CycleYear,
			Description: null.StringFrom("Premium yearly price"),
		}, Tier: enum.TierPremium,
		IsActive:   true,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  "weiguo.ni",
	},
	Discount: Discount{
		DiscID: null.String{},
		DiscountInput: DiscountInput{
			PriceOff: null.Float{},
			Percent:  null.Int{},
			Period: Period{
				StartUTC: chrono.Time{},
				EndUTC:   chrono.Time{},
			},
		},
	},
}

func mustStringify(v interface{}) []byte {
	b, err := json.MarshalIndent(v, "", "\t")

	if err != nil {
		panic(err)
	}

	return b
}

func seedGoFake() {
	gofakeit.Seed(time.Now().UnixNano())
}

func TestNewBanner(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())

	input := BannerInput{
		Heading:    gofakeit.Sentence(10),
		CoverURL:   null.StringFrom(gofakeit.URL()),
		SubHeading: null.StringFrom(gofakeit.Sentence(5)),
		Content:    null.StringFrom(gofakeit.Paragraph(3, 2, 5, "\n")),
	}

	banner := NewBanner(input, gofakeit.Username())

	assert.NotEmpty(t, banner.Heading)

	t.Logf("Request data: %s", mustStringify(input))
}

func TestNewPromo(t *testing.T) {
	gofakeit.Seed(time.Now().UnixNano())
	input := PromoInput{
		Heading:    null.StringFrom(gofakeit.Sentence(10)),
		CoverURL:   null.StringFrom(gofakeit.URL()),
		SubHeading: null.StringFrom(gofakeit.Sentence(5)),
		Content:    null.StringFrom(gofakeit.Paragraph(3, 2, 5, "\n")),
		Period: Period{
			StartUTC: chrono.TimeNow(),
			EndUTC:   chrono.TimeNow(),
		},
	}

	p := NewPromo(input, gofakeit.Username())

	assert.NotEmpty(t, p.ID)

	t.Logf("Promo input %s", mustStringify(input))
}
