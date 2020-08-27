package products

import (
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			args:    args{b: test.NewPaywallBanner()},
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

func TestEnv_UpdateBanner(t *testing.T) {
	b := test.NewPaywallBanner()
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

func TestEnv_CreatePromo(t *testing.T) {

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		bannerID int64
		p        paywall.Promo
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
				bannerID: 1,
				p:        test.NewPaywallPromo(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.CreatePromo(tt.args.bannerID, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("CreatePromo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadPromo(t *testing.T) {
	p := test.NewPaywallPromo()

	test.NewRepo().MustCreatePromo(p)

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
		wantErr bool
	}{
		{
			name:    "Load promo",
			fields:  fields{db: test.DBX},
			args:    args{id: p.ID},
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

			t.Logf("Retrieved promo: %+v", p)

			assert.NotEmpty(t, got.ID)
		})
	}
}

func TestEnv_DropBannerPromo(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		bannerID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Drop promo from banner",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				bannerID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.DropBannerPromo(tt.args.bannerID); (err != nil) != tt.wantErr {
				t.Errorf("DropBannerPromo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_listPaywallProducts(t *testing.T) {

	_ = test.NewRepo().CreatePaywallProducts()

	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Retrieve paywall products",
			fields: fields{
				db: test.DBX,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.retrievePaywallProducts()
			if (err != nil) != tt.wantErr {
				t.Errorf("retrievePaywallProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Len(t, got, 2)
		})
	}
}

func TestEnv_listPaywallPlans(t *testing.T) {
	_ = test.NewRepo().CreatePaywallProducts()

	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Retrieve all plans",
			fields: fields{
				db: test.DBX,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.retrievePaywallPlans()
			if (err != nil) != tt.wantErr {
				t.Errorf("retrievePaywallPlans() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Len(t, got, 3)
		})
	}
}

func TestEnv_LoadPaywallProducts(t *testing.T) {

	_ = test.NewRepo().CreatePaywallProducts()

	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Load paywall products with prices",
			fields: fields{
				db: test.DBX,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.LoadPaywallProducts()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPaywallProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Len(t, got, 2)
		})
	}
}
