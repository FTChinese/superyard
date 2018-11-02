package subscription

import (
	"database/sql"
	"testing"
)

func newDevEnv() Env {
	db, err := sql.Open("mysql", "sampadm:secret@unix(/tmp/mysql.sock)/")

	if err != nil {
		panic(err)
	}

	return Env{DB: db}
}

var devEnv = newDevEnv()

func TestNewDiscount(t *testing.T) {
	d := Discount{
		Name:        "Double Eleven",
		Description: "15 discount",
		Start:       "2018-11-10T16:00:00Z",
		End:         "2018-11-11T16:00:00Z",
		Plans: map[string]Plan{
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
		},
	}

	err := devEnv.NewDiscount(d)

	if err != nil {
		t.Error(err)
	}
}

func TestRetrieveDiscount(t *testing.T) {
	d, err := devEnv.Retrieve(2)

	if err != nil {
		t.Error(err)
	}

	t.Log(d)
}
