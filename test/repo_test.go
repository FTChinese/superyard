package test

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/superyard/pkg/subs"
	"testing"
)

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
			args:    args{order: NewPersona().Order(true)},
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
		})
	}
}
