package user

import (
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/staff"
	"github.com/FTChinese/superyard/test"
	"testing"
)

func TestEnv_VerifyPassword(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DBs: db.MustNewMyDBs(false)}

	type args struct {
		verifier staff.PasswordVerifier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Verify password",
			args: args{
				verifier: staff.PasswordVerifier{
					StaffID:     s.ID,
					OldPassword: s.Password,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.VerifyPassword(tt.args.verifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Account: %+v", got)
		})
	}
}

func TestEnv_UpdatePassword(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DBs: db.MustNewMyDBs(false)}

	type args struct {
		c staff.Credentials
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Update password",
			args:    args{c: s.Credentials()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdatePassword(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
