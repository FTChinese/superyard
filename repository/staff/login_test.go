package staff

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/employee"
	"gitlab.com/ftchinese/backyard-api/test"
	"testing"
)

func TestEnv_Login(t *testing.T) {

	s := test.NewStaff()

	createStaff(s.Account())

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
		wantErr bool
	}{
		{
			name: "Staff Login",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				l: s.Login(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.Login(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Got: %+v", got)
		})
	}
}
