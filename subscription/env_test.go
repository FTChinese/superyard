package subscription

import (
	"database/sql"

	"github.com/icrowley/fake"
)

func newDevEnv() Env {
	db, err := sql.Open("mysql", "sampadm:secret@unix(/tmp/mysql.sock)/")

	if err != nil {
		panic(err)
	}

	return Env{DB: db}
}

var mockSchedule = Schedule{
	Name:        fake.Brand(),
	Description: fake.Product(),
	Start:       "2018-11-07T16:00:00Z",
	End:         "2018-11-11T16:00:00Z",
}

var mockPricing = map[string]Plan{
	"standard_year": Plan{
		Tier:        "standard",
		Cycle:       "year",
		Price:       0.01,
		ID:          10,
		Description: "FT中文网 - 标准会员",
		Message:     "Double Eleben Discount",
	},
	"standard_month": Plan{
		Tier:        "standard",
		Cycle:       "month",
		Price:       0.01,
		ID:          5,
		Description: "FT中文网 - 标准会员",
		Message:     "Double Eleben Discount",
	},
	"premium_year": Plan{
		Tier:        "premium",
		Cycle:       "year",
		Price:       0.01,
		ID:          100,
		Description: "FT中文网 - 高端会员",
		Message:     "Double Eleben Discount",
	},
}

var mockBanner = Banner{
	Heading:    "FT中文网会员订阅服务",
	SubHeading: "欢迎您",
	Content: []string{
		"希望全球视野的FT中文网，能够带您站在高海拔的地方俯瞰世界，引发您的思考，从不同的角度看到不一样的事物，见他人之未见！",
	},
}

var devEnv = newDevEnv()
