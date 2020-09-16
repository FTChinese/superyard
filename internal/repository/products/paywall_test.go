package products

import (
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_retrievePaywallPromo(t *testing.T) {
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
			name: "Get paywall promo",
			fields: fields{
				db: test.DBX,
			},
			args:    args{bannerID: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.retrievePaywallPromo(tt.args.bannerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("retrievePaywallPromo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotEmpty(t, got.ID)
		})
	}
}

func TestEnv_retrievePaywallProducts(t *testing.T) {

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

func TestEnv_retrievePaywallPlans(t *testing.T) {
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

func TestEnv_LoadPaywall(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Load paywall",
			fields: fields{
				db: test.DBX,
			},
			args:    args{id: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.LoadPaywall(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPaywall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}