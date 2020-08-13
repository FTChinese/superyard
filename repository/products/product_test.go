package products

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestEnv_CreatePricedProduct(t *testing.T) {

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		p paywall.PricedProduct
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create a standard with optional prices",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				p: test.NewProductMocker(enum.TierStandard).PricedProduct(),
			},
			wantErr: false,
		},
		{
			name: "Create a premium with optional prices",
			fields: fields{
				db: test.DBX,
			},
			args: args{
				p: test.NewProductMocker(enum.TierPremium).PricedProduct(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			env := Env{
				db: tt.fields.db,
			}
			if err := env.CreatePricedProduct(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("CreatePricedProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadProduct(t *testing.T) {
	prod := test.NewProductMocker(enum.TierStandard).Product()
	_ = test.NewRepo().CreateProduct(prod)

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
			name:    "Load product",
			fields:  fields{db: test.DBX},
			args:    args{id: prod.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.LoadProduct(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.ID, prod.ID)
		})
	}
}

func TestEnv_UpdateProduct(t *testing.T) {
	prod := test.NewProductMocker(enum.TierStandard).Product()
	_ = test.NewRepo().CreateProduct(prod)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		prod paywall.Product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Update product",
			fields: fields{db: test.DBX},
			args: args{
				prod: prod.Update(
					test.NewProductMocker(enum.TierStandard).
						Product().
						ProductInput,
				),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.UpdateProduct(tt.args.prod); (err != nil) != tt.wantErr {
				t.Errorf("UpdateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_ActivateProduct(t *testing.T) {
	prodStd := test.NewProductMocker(enum.TierStandard).Product()
	prodPrm := test.NewProductMocker(enum.TierPremium).Product()

	repo := test.NewRepo()
	_ = repo.CreateProduct(prodStd)
	_ = repo.CreateProduct(prodPrm)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		prod paywall.Product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Activate standard",
			fields:  fields{db: test.DBX},
			args:    args{prod: prodStd},
			wantErr: false,
		},
		{
			name:    "Activate premium",
			fields:  fields{db: test.DBX},
			args:    args{prod: prodPrm},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			if err := env.ActivateProduct(tt.args.prod); (err != nil) != tt.wantErr {
				t.Errorf("ActivateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_ListPricedProducts(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "List product with pricing plans",
			fields:  fields{db: test.DBX},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				db: tt.fields.db,
			}
			got, err := env.ListPricedProducts()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPricedProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}
