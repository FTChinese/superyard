package staff

import (
	"gitlab.com/ftchinese/backyard-api/models/employee"
	"gitlab.com/ftchinese/backyard-api/test"
	"testing"
)

func TestEnv_Activate(t *testing.T) {
	env := Env{DB: test.DBX}

	a := test.NewStaff().Account()
	if err := env.Create(a); err != nil {
		t.Error(err)
	}
	if err := env.Deactivate(a.ID.String); err != nil {
		t.Error(err)
	}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Activate",
			args:    args{id: a.ID.String},
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

func TestEnv_AddID(t *testing.T) {
	env := Env{DB: test.DBX}

	a := test.NewStaff().Account()
	if err := env.Create(a); err != nil {
		t.Error(err)
	}

	type args struct {
		a employee.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Add id",
			args:    args{a: a},
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

func TestEnv_Create(t *testing.T) {

	env := Env{DB: test.DBX}

	type args struct {
		a employee.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create a new staff",
			args:    args{a: test.NewStaff().Account()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.Create(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_Deactivate(t *testing.T) {
	env := Env{DB: test.DBX}

	a := test.NewStaff().Account()
	if err := env.Create(a); err != nil {
		t.Error(err)
	}

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Deactivate",
			args:    args{id: a.ID.String},
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
