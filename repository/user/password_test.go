package user

import (
	"gitlab.com/ftchinese/superyard/models/staff"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_UpdatePassword(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

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
			args:    args{c: s.NewPassword()},
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

func TestEnv_VerifyPassword(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

	type args struct {
		c staff.Credentials
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Verify password",
			args:    args{c: s.OldPassword()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.VerifyPassword(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Account: %+v", got)
		})
	}
}
