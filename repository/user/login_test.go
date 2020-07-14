package user

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/superyard/pkg/staff"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_Login(t *testing.T) {

	s := test.NewStaff()

	repo := test.NewRepo()
	repo.MustCreateStaff(s.SignUp())

	env := Env{DB: test.DBX}

	type args struct {
		l staff.Credentials
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Login",
			args:    args{l: s.Credentials()},
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

	repo := test.NewRepo()

	s := test.NewStaff()
	repo.MustCreateStaff(s.SignUp())

	env := Env{DB: test.DBX}

	type args struct {
		l  staff.Credentials
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
				l:  s.Credentials(),
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

func TestEnv_SavePwResetSession(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DB: test.DBX}

	type args struct {
		session staff.PwResetSession
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create a new password reset session",
			args: args{
				session: s.PwResetSession(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SavePwResetSession(tt.args.session); (err != nil) != tt.wantErr {
				t.Errorf("SavePwResetSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadPwResetSession(t *testing.T) {
	s := test.NewStaff()

	test.NewRepo().MustSavePwResetSession(s.PwResetSession())

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		token string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    staff.PwResetSession
		wantErr bool
	}{
		{
			name: "Load a password reset session",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				token: s.PwResetToken,
			},
			//want:    staff.PwResetSession{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.LoadPwResetSession(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPwResetSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoadPwResetSession() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%+v", got)
		})
	}
}

func TestEnv_AccountByResetToken(t *testing.T) {
	s := test.NewStaff()

	repo := test.NewRepo()
	repo.MustCreateStaff(s.SignUp())
	repo.MustSavePwResetSession(s.PwResetSession())

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		token string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		//want    staff.Account
		wantErr bool
	}{
		{
			name: "Load account for a password reset token",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				token: s.PwResetToken,
			},
			//want:    staff.Account{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.AccountByResetToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountByResetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("AccountByResetToken() got = %v, want %v", got, tt.want)
			//}

			t.Logf("%+v", got)
		})
	}
}

func TestEnv_DisableResetToken(t *testing.T) {

	s := test.NewStaff()

	repo := test.NewRepo()
	repo.MustCreateStaff(s.SignUp())
	repo.MustSavePwResetSession(s.PwResetSession())

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Disable a password reset token",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				token: s.PwResetToken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			if err := env.DisableResetToken(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("DisableResetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
