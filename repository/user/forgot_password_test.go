package user

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_SavePwResetToken(t *testing.T) {
	env := Env{DB: test.DBX}
	s := test.NewStaff()

	test.NewRepo().MustCreateStaff(s)

	type args struct {
		pr employee.PasswordReset
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Save email + token",
			args: args{pr: s.PasswordReset()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.SavePwResetToken(tt.args.pr); (err != nil) != tt.wantErr {
				t.Errorf("SavePwResetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_AccountByResetToken(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

	_ = env.SavePwResetToken(s.PasswordReset())

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Get account for a password reset token",
			args:    args{token: s.PwResetToken},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.AccountByResetToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountByResetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Account: %+v", got)
		})
	}
}

func TestEnv_DisableResetToken(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

	_ = env.SavePwResetToken(s.PasswordReset())

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Disable a password reset token",
			args:    args{token: s.PwResetToken},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.DisableResetToken(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("DisableResetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
