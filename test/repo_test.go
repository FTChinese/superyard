package test

import (
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestRepo_CreateReader(t *testing.T) {

	p := NewPersona()

	t.Logf("Creating user %s", p.FtcID)

	err := NewRepo().CreateReader(p.FtcAccount())
	if err != nil {
		t.Error(err)
	}
}

func TestRepo_CreateVIP(t *testing.T) {
	repo := NewRepo()

	err := repo.CreateVIP(NewPersona().SetVIP().FtcAccount())
	if err != nil {
		t.Error(err)
	}
}

func TestRepo_CreateOrder(t *testing.T) {
	order := NewPersona().Order(false)

	err := NewRepo().CreateOrder(order)
	if err != nil {
		t.Error(err)
	}
}

func TestRepo_CreatePaywallProducts(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Create all product data",
			fields: fields{
				db: DBX,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repo{
				db: tt.fields.db,
			}
			if err := repo.CreatePaywallProducts(); (err != nil) != tt.wantErr {
				t.Errorf("CreatePaywallProducts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
