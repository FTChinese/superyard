package subsapi

import (
	"bytes"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/pkg/dt"
	"github.com/FTChinese/superyard/pkg/price"
	"github.com/FTChinese/superyard/test"
	"github.com/guregu/null"
	"io"
	"testing"
	"time"
)

func TestClient_RefreshPaywall(t *testing.T) {
	faker.MustConfigViper()

	c := NewClient(false)

	resp, err := c.RefreshPaywall()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Status %d", resp.StatusCode)
	t.Logf("%s", faker.MustReadBody(resp.Body))
}

func TestClient_LoadPaywall(t *testing.T) {
	faker.MustConfigViper()

	c := NewClient(false)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Load paywall",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.LoadPaywall()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPaywall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Resonse status %d", got.StatusCode)
			t.Logf("%s", faker.MustReadBody(got.Body))
		})
	}
}

func TestClient_CreatePrice(t *testing.T) {
	faker.MustConfigViper()

	c := NewClient(false)

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create price",
			args: args{
				body: test.NewProductBuilder("").
					NewPriceBuilder("").
					BuildIOBody(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.CreatePrice(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Resonse status %d", got.StatusCode)
			t.Logf("%s", faker.MustReadBody(got.Body))
		})
	}
}

func TestClient_ActivatePrice(t *testing.T) {
	faker.MustConfigViper()

	c := NewClient(false)

	type args struct {
		priceID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Activate price",
			args: args{
				priceID: "plan_RKy1IuKSXyua",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.ActivatePrice(tt.args.priceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ActivatePrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Resonse status %d", got.StatusCode)
			t.Logf("%s", faker.MustReadBody(got.Body))
		})
	}
}

func TestClient_RefreshPriceDiscounts(t *testing.T) {

	faker.MustConfigViper()

	c := NewClient(false)
	type args struct {
		priceID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Refresh price discounts",
			args: args{
				priceID: "plan_rLIy6LJYW8LV",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.RefreshPriceDiscounts(tt.args.priceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshPriceDiscounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Resonse status %d", got.StatusCode)
			t.Logf("%s", faker.MustReadBody(got.Body))
		})
	}
}

func TestClient_CreateDiscount(t *testing.T) {
	faker.MustConfigViper()

	c := NewClient(false)

	type args struct {
		body io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create retention",
			args: args{
				body: bytes.NewReader(faker.MustMarshalIndent(
					price.DiscountParams{
						CreatedBy:      "weiguo.ni",
						Description:    null.StringFrom("现在续订享75折优惠"),
						Kind:           price.OfferKindRetention,
						Percent:        null.IntFrom(75),
						DateTimePeriod: dt.DateTimePeriod{},
						PriceOff:       null.FloatFrom(80),
						PriceID:        "plan_RKy1IuKSXyua",
						Recurring:      false,
					},
				)),
			},
			wantErr: false,
		},
		{
			name: "Create win back",
			args: args{
				body: bytes.NewReader(faker.MustMarshalIndent(
					price.DiscountParams{
						CreatedBy:      "weiguo.ni",
						Description:    null.StringFrom("现现在购买享85折优惠"),
						Kind:           price.OfferKindWinBack,
						Percent:        null.IntFrom(85),
						DateTimePeriod: dt.DateTimePeriod{},
						PriceOff:       null.FloatFrom(40),
						PriceID:        "plan_RKy1IuKSXyua",
						Recurring:      false,
					},
				)),
			},
			wantErr: false,
		},
		{
			name: "Create Introductory",
			args: args{
				body: bytes.NewReader(faker.MustMarshalIndent(
					price.DiscountParams{
						CreatedBy:   "weiguo.ni",
						Description: null.StringFrom("新会员订阅仅需1元"),
						Kind:        price.OfferKindIntroductory,
						Percent:     null.Int{},
						DateTimePeriod: dt.DateTimePeriod{
							StartUTC: chrono.TimeNow(),
							EndUTC:   chrono.TimeFrom(time.Now().AddDate(0, 0, 7)),
						},
						PriceOff:  null.FloatFrom(34),
						PriceID:   "plan_ohky3lyEMPSf",
						Recurring: false,
					},
				)),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := c.CreateDiscount(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDiscount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Resonse status %d", got.StatusCode)
			t.Logf("%s", faker.MustReadBody(got.Body))
		})
	}
}

func TestClient_RemoveDiscount(t *testing.T) {
	faker.MustConfigViper()

	c := NewClient(false)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Remove a discount",
			args: args{
				id: "dsc_7KiZow3Hj7G1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.RemoveDiscount(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveDiscount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Resonse status %d", got.StatusCode)
			t.Logf("%s", faker.MustReadBody(got.Body))
		})
	}
}
