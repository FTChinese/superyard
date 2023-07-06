package auth

import (
	"reflect"
	"testing"

	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/db"
)

func TestEnv_VerifyPassword(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a := mockNewAccount()
	a, _ = env.createdAccount(a)

	type args struct {
		id          int64
		currentPass string
	}
	tests := []struct {
		name    string
		args    args
		want    user.Account
		wantErr bool
	}{
		{
			name: "verify password",
			args: args{
				id:          a.ID,
				currentPass: a.Password,
			},
			want:    a.Account,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.VerifyPassword(tt.args.id, tt.args.currentPass)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_UpdatePassword(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	a := mockNewAccount()
	env.createdAccount(a)

	type args struct {
		holder user.Credentials
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Update password",
			args: args{
				holder: user.Credentials{
					UserName: a.UserName,
					Password: faker.SimplePassword(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			t.Logf("Change password from %s to %s", tt.args.holder.UserName, tt.args.holder.Password)

			if err := env.UpdatePassword(tt.args.holder); (err != nil) != tt.wantErr {
				t.Errorf("Env.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
