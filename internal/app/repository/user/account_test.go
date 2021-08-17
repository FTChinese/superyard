package user

import (
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/staff"
	"github.com/FTChinese/superyard/test"
	"testing"
)

func TestEnv_AccountByID(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DBs: db.MustNewMyDBs(false)}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Get user by staff id",
			args:    args{id: s.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.AccountByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Accout: %+v", got)
		})
	}
}

func TestEnv_AccountByEmail(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DBs: db.MustNewMyDBs(false)}

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Get account by name",
			args:    args{email: s.Email},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := env.AccountByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Account: %+v", got)
		})
	}
}

func TestEnv_AddID(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DBs: db.MustNewMyDBs(false)}

	type args struct {
		a staff.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Add staff id",
			args:    args{a: s.Account()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.AddID(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("AddID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_SetEmail(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DBs: db.MustNewMyDBs(false)}

	type args struct {
		a staff.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Set email",
			args:    args{a: s.Account()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.SetEmail(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("SetEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_UpdateDisplayName(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DBs: db.MustNewMyDBs(false)}

	type args struct {
		a staff.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Update display name",
			args:    args{a: s.Account()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateDisplayName(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDisplayName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_RetrieveProfile(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s.SignUp())

	env := Env{DBs: db.MustNewMyDBs(false)}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve profile",
			args:    args{id: s.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveProfile(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Profile: %+v", got)
		})
	}
}
