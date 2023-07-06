package auth

import (
	"reflect"
	"testing"

	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/db"
)

func mockCredentials(a mockAccount) user.Credentials {
	return user.Credentials{
		UserName: a.UserName,
		Password: a.Password,
	}
}

func TestEnv_Login(t *testing.T) {
	env := NewEnv(db.MockGormSQL())
	a := mockNewAccount()
	c := mockCredentials(a)
	env.createdAccount(a)

	type args struct {
		c user.Credentials
	}
	tests := []struct {
		name    string
		args    args
		want    user.Account
		wantErr bool
	}{
		{
			name: "mock login",
			args: args{
				c: c,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.Login(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_SavePwResetSession(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	type args struct {
		session user.PwResetSession
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "save password reset token",
			args: args{
				session: user.MustNewPwResetSession("neefrankie@163.com"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SavePwResetSession(tt.args.session); (err != nil) != tt.wantErr {
				t.Errorf("Env.SavePwResetSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_LoadPwResetSession(t *testing.T) {

	env := NewEnv(db.MockGormSQL())

	sess := user.MustNewPwResetSession("neefrankie@163.com")

	_ = env.SavePwResetSession(sess)

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    user.PwResetSession
		wantErr bool
	}{
		{
			name: "Load a reset session",
			args: args{
				token: sess.Token.String(),
			},
			want:    sess,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.LoadPwResetSession(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.LoadPwResetSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.LoadPwResetSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_DisableResetToken(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	sess := user.MustNewPwResetSession("neefrankie@163.com")

	_ = env.SavePwResetSession(sess)

	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "disabled reset token",
			args: args{
				token: sess.Token.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.DisableResetToken(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("Env.DisableResetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
