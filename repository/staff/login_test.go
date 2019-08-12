package staff

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/employee"
	"gitlab.com/ftchinese/backyard-api/test"
	"testing"
)

func TestEnv_VerifyPassword(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		l employee.Login
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:   "Login",
			fields: fields{DB: test.DBX},
			args: args{l: employee.Login{
				UserName: "nobis",
				Password: "12345678",
			}},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.VerifyPassword(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
