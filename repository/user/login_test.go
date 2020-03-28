package user

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_Login(t *testing.T) {
	env := Env{DB: test.DBX}

	repo := test.NewRepo()
	s := test.NewStaff()
	repo.MustCreateStaff(s)

	type args struct {
		l employee.Login
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Login",
			args:    args{l: s.Login()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.Login(tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Account: %+v", got)
		})
	}
}

func TestEnv_UpdateLastLogin(t *testing.T) {
	env := Env{DB: test.DBX}

	repo := test.NewRepo()
	s := test.NewStaff()
	repo.MustCreateStaff(s)

	type args struct {
		l  employee.Login
		ip string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Login IP",
			args: args{
				l:  s.Login(),
				ip: s.IP,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.UpdateLastLogin(tt.args.l, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("UpdateLastLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}