package test

import (
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/jmoiron/sqlx"
	"testing"
)

func TestRepo_CreateReader(t *testing.T) {

	p := NewPersona()

	t.Logf("Creating user %s", p.FtcID)

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		u reader.SandboxUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Creat a user",
			fields:  fields{db: DBX},
			args:    args{u: p.Reader()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repo{
				db: tt.fields.db,
			}
			if err := repo.CreateReader(tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("CreateReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_CreateOrder(t *testing.T) {
	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		order subs.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create order",
			fields: fields{
				db: DBX,
			},
			args:    args{order: NewPersona().Order(false)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repo{
				db: tt.fields.db,
			}
			if err := repo.CreateOrder(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}

			t.Logf("New order id: %s", tt.args.order.ID)
		})
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
