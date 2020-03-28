package admin

import (
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_AccountByID(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Get account by id",
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

			t.Logf("Account: %+v", got)
		})
	}
}

func TestEnv_UpdateAccount(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

	type args struct {
		p employee.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Update account",
			args:    args{p: s.Account()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateAccount(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_Deactivate(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Deactivate an account",
			args:    args{id: s.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.Deactivate(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Deactivate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_Activate(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Activate an account",
			args:    args{id: s.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.Activate(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Activate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_AccountByName(t *testing.T) {
	s := test.NewStaff()
	test.NewRepo().MustCreateStaff(s)

	env := Env{DB: test.DBX}

	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Get account by name",
			args: args{name: s.UserName},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.AccountByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Account: %+v", got)
		})
	}
}