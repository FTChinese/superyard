package subscription

import (
	"database/sql"
	"time"

	"github.com/bxcodec/faker"

	"github.com/icrowley/fake"
	"gitlab.com/ftchinese/backyard-api/util"
)

func newDevEnv() Env {
	db, err := sql.Open("mysql", "sampadm:secret@unix(/tmp/mysql.sock)/")

	if err != nil {
		panic(err)
	}

	return Env{DB: db}
}

var tomrrow = time.Now().AddDate(0, 0, 1)

var mockSchedule = Schedule{
	Name:        fake.Brand(),
	Description: fake.Product(),
	Start:       util.ISO8601UTC.FromTime(tomrrow),
	End:         util.ISO8601UTC.FromTime(tomrrow.AddDate(0, 0, 1)),
}

var mockPricing = map[string]Plan{
	keyStdYear: Plan{
		Tier:        TierStandard,
		Cycle:       CycleYear,
		Price:       198,
		ID:          10,
		Description: "FT中文网 - 标准会员",
		Message:     "",
		Ignore:      false,
	},
	keyStdMonth: Plan{
		Tier:        TierStandard,
		Cycle:       CycleMonth,
		Price:       28,
		ID:          5,
		Description: "FT中文网 - 标准会员",
		Message:     "",
		Ignore:      true,
	},
	keyPrmYear: Plan{
		Tier:        TierPremium,
		Cycle:       CycleYear,
		Price:       1998,
		ID:          100,
		Description: "FT中文网 - 高端会员",
		Message:     "",
		Ignore:      false,
	},
}

var mockBanner = Banner{
	CoverURL:   faker.Internet{}.URL(),
	Heading:    fake.Sentence(),
	SubHeading: fake.Sentence(),
	Content: []string{
		fake.Paragraph(),
		fake.Paragraph(),
	},
}

var devEnv = newDevEnv()

func createPromo() (int64, error) {
	id, err := devEnv.NewSchedule(mockSchedule, "weiguo.ni")

	if err != nil {
		return id, err
	}

	err = devEnv.SavePricing(id, mockPricing)

	if err != nil {
		return id, err
	}

	err = devEnv.SaveBanner(id, mockBanner)

	if err != nil {
		return id, err
	}

	return id, nil
}
