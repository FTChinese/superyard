package auth

import (
	"reflect"
	"testing"

	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/db"
)

func TestEnv_SavePwResetSession(t *testing.T) {
	faker.MustSetupViper()
	env := Env{
		gormDBs: db.MustNewMultiGormDBs(false),
	}

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
	faker.MustSetupViper()
	env := Env{
		gormDBs: db.MustNewMultiGormDBs(false),
	}

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
				token: string(sess.Token),
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
