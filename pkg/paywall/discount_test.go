package paywall

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"testing"
	"time"
)

func TestNewDiscountSchema(t *testing.T) {
	seedGoFake()

	input := DiscountInput{
		PriceOff: null.FloatFrom(59),
		Percent:  null.Int{},
		Period: Period{
			StartUTC: chrono.TimeNow(),
			EndUTC:   chrono.TimeFrom(time.Now().AddDate(0, 0, 1)),
		},
	}

	schema := NewDiscountSchema(input, genPlanID(), gofakeit.Username())

	t.Logf("%s", mustStringify(schema))
}
