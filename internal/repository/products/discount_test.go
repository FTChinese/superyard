package products

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestEnv_CreateDiscount(t *testing.T) {

	pm := test.NewProductMocker(enum.TierStandard)
	plan := pm.Plan(enum.CycleYear)

	repo := test.NewRepo()
	_ = repo.CreateProduct(pm.Product())
	_ = repo.CreatePlan(plan)

	env := NewEnv(db.MustNewMyDBs(false))

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		d paywall.DiscountSchema
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create a discount",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				d: test.NewDiscount(plan),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.CreateDiscount(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("CreateDiscount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_DropDiscount(t *testing.T) {
	pm := test.NewProductMocker(enum.TierStandard)
	plan := pm.Plan(enum.CycleYear)

	t.Logf("Plan id %s", plan.ID)

	repo := test.NewRepo()
	_ = repo.CreateProduct(pm.Product())
	_ = repo.CreatePlan(plan)

	env := NewEnv(db.MustNewMyDBs(false))
	env.CreateDiscount(test.NewDiscount(plan))

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		plan paywall.Plan
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{

		{
			name: "Remove discount from plan",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				plan: plan,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.DropDiscount(tt.args.plan); (err != nil) != tt.wantErr {
				t.Errorf("DropDiscount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
