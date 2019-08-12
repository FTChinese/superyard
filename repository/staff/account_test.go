package staff

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/employee"
	"gitlab.com/ftchinese/backyard-api/test"
	"testing"
)

func TestEnv_Create(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		a employee.Account
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Create account",
			fields:  fields{DB: test.DBX},
			args:    args{a: test.GenEmployee()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			if err := env.Create(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
