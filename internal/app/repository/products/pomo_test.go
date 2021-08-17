package products

import (
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_CreatePromo(t *testing.T) {
	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		bannerID int64
		p        paywall.Promo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create promo",
			args: args{
				bannerID: 1,
				p:        test.NewPaywallPromo(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.CreatePromo(tt.args.bannerID, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("CreatePromo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadPromo(t *testing.T) {
	p := test.NewPaywallPromo()

	test.NewRepo().MustCreatePromo(p)
	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Load promo",
			args:    args{id: p.ID.String},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		bannerID int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Drop promo from banner",
			args: args{
				bannerID: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.DropBannerPromo(tt.args.bannerID); (err != nil) != tt.wantErr {
				t.Errorf("DropBannerPromo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}