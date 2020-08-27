package products

import (
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
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
