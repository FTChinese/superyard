package products

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/test"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func newBannerInput() paywall.BannerInput {
	test.SeedGoFake()

	return paywall.BannerInput{
		Heading:    gofakeit.Sentence(10),
		CoverURL:   null.StringFrom(gofakeit.URL()),
		SubHeading: null.StringFrom(gofakeit.Sentence(5)),
		Content:    null.StringFrom(gofakeit.Paragraph(3, 2, 5, "\n")),
	}
}

func newPeriod() paywall.Period {
	return paywall.Period{
		StartUTC: chrono.TimeNow(),
		EndUTC:   chrono.TimeFrom(time.Now().AddDate(0, 0, 1)),
	}
}

func newPromoInput() paywall.PromoInput {
	test.SeedGoFake()

	return paywall.PromoInput{
		BannerInput: newBannerInput(),
		Period:      newPeriod(),
	}
}

func TestEnv_CreateBanner(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		b paywall.Banner
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Generate banner",
			fields: fields{
				db: test.DBX,
			},
			args:    args{b: paywall.NewBanner(newBannerInput(), "weiguo.ni")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.CreateBanner(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("CreateBanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadBanner(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    paywall.Banner
		wantErr bool
	}{
		{
			name:    "Load banner",
			fields:  fields{db: test.DBX},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.LoadBanner(1)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Banner: %+v", got)
		})
	}
}

func TestEnv_UpdateBanner(t *testing.T) {
	b := paywall.NewBanner(newBannerInput(), "weiguo.ni")
	b.ID = 1

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		b paywall.Banner
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update banner",
			fields: fields{
				db: test.DBX,
			},
			args:    args{b: b},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.UpdateBanner(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_CreatePromo(t *testing.T) {
	p := paywall.NewPromo(newPromoInput(), "weiguo.ni")

	t.Logf("%+v", p)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		p paywall.Promo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Create promo",
			fields: fields{db: test.DBX},
			args: args{
				p: p,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.CreatePromo(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("CreatePromo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadPromo(t *testing.T) {
	env := NewEnv(test.DBX)
	promo := paywall.NewPromo(newPromoInput(), "weiguo.ni")

	_ = env.CreatePromo(promo)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    paywall.Promo
		wantErr bool
	}{
		{
			name:    "Load promo",
			fields:  fields{db: test.DBX},
			args:    args{id: promo.ID},
			want:    promo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.LoadPromo(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPromo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.ID, tt.want.ID)
		})
	}
}
