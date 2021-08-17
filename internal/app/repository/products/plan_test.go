package products

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_CreatePlan(t *testing.T) {
	pm := test.NewProductMocker(enum.TierStandard)
	_ = test.NewRepo().CreateProduct(pm.Product())

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		p paywall.Plan
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create plan",
			args:    args{p: pm.Plan(enum.CycleYear)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.CreatePlan(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("CreatePlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadPlan(t *testing.T) {
	pm := test.NewProductMocker(enum.TierStandard)
	plan := pm.Plan(enum.CycleYear)

	repo := test.NewRepo()
	_ = repo.CreateProduct(pm.Product())
	_ = repo.CreatePlan(plan)

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
			name: "Load plan",
			args: args{
				id: plan.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadPlan(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.ID, plan.ID)
		})
	}
}

func TestEnv_ActivatePlan(t *testing.T) {
	pm := test.NewProductMocker(enum.TierStandard)
	plan1 := pm.Plan(enum.CycleYear)
	plan2 := pm.Plan(enum.CycleMonth)
	plan3 := pm.Plan(enum.CycleMonth)

	repo := test.NewRepo()
	_ = repo.CreateProduct(pm.Product())
	_ = repo.CreatePlan(plan1)
	_ = repo.CreatePlan(plan2)
	_ = repo.CreatePlan(plan3)

	env := NewEnv(db.MustNewMyDBs(false))

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
			name: "Activate yearly plan",
			args: args{
				plan: plan1,
			},
			wantErr: false,
		},
		{
			name: "Activate monthly plan",
			args: args{
				plan: plan2,
			},
			wantErr: false,
		},
		{
			name: "Modify previous one",
			args: args{
				plan: plan3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("Activate plan %s", tt.args.plan.ID)
			if err := env.ActivatePlan(tt.args.plan); (err != nil) != tt.wantErr {
				t.Errorf("ActivatePlan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

func TestEnv_ListPlansOfProduct(t *testing.T) {
	pm := test.NewProductMocker(enum.TierStandard)
	prod := pm.Product()
	plan1 := pm.Plan(enum.CycleYear)
	plan2 := pm.Plan(enum.CycleMonth)

	repo := test.NewRepo()
	_ = repo.CreateProduct(prod)
	_ = repo.CreateAndActivatePlan(plan1)
	_ = repo.CreateAndActivatePlan(plan2)

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		productID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "List plans of a product",
			args: args{
				productID: prod.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListPlansOfProduct(tt.args.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPlansOfProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Len(t, got, 2)

			assert.NotEmpty(t, got[0].ID)
		})
	}
}
