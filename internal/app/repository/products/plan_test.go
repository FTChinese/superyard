package products

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/test"
	"testing"
)

func TestEnv_ProductHasPlan(t *testing.T) {
	repo := test.NewRepo()

	pm1 := test.NewProductMocker(enum.TierStandard)
	prod1 := pm1.Product()
	_ = repo.CreateProduct(prod1)
	_ = repo.CreateAndActivatePlan(pm1.Plan(enum.CycleYear))

	pm2 := test.NewProductMocker(enum.TierPremium)
	prod2 := pm2.Product()
	_ = repo.CreateProduct(prod2)

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		productID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Product has active plan",
			args: args{
				productID: prod1.ID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Product has no active plan",
			args: args{
				productID: prod2.ID,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ProductHasActivePlan(tt.args.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductHasActivePlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProductHasActivePlan() got = %v, want %v", got, tt.want)
			}
		})
	}
}
