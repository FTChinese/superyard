package model

import (
	"database/sql"
	"log"
	"testing"

	"github.com/icrowley/fake"
	"gitlab.com/ftchinese/backyard-api/staff"
)

func TestStaffEnv_UpdateLoginHistory(t *testing.T) {
	mock := newMockStaff()
	mock.createAccount()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		l  staff.Login
		ip string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Update Login History",
			fields:  fields{DB: db},
			args:    args{l: mock.login(), ip: fake.IPv4()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			if err := env.UpdateLoginHistory(tt.args.l, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.UpdateLoginHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffEnv_LoadAccountByName(t *testing.T) {
	mock := newMockStaff()
	mock.createAccount()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		name   string
		active bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Load Account by Name",
			fields:  fields{DB: db},
			args:    args{name: mock.userName, active: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			got, err := env.LoadAccountByName(tt.args.name, tt.args.active)
			if (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.LoadAccountByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			log.Printf("%+v", got)
		})
	}
}

func TestStaffEnv_UpdateName(t *testing.T) {
	mock := newMockStaff()
	mock.createAccount()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName    string
		displayName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Update Name",
			fields: fields{DB: db},
			args: args{
				userName:    mock.userName,
				displayName: fake.FullName(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			if err := env.UpdateName(tt.args.userName, tt.args.displayName); (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.UpdateName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffEnv_UpdateEmail(t *testing.T) {
	mock := newMockStaff()
	mock.createAccount()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName string
		email    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Update Email",
			fields:  fields{DB: db},
			args:    args{userName: mock.userName, email: fake.EmailAddress()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			if err := env.UpdateEmail(tt.args.userName, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.UpdateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffEnv_Profile(t *testing.T) {
	mocker := newMockStaff()
	mocker.createAccount()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Load Staff Profile",
			fields:  fields{DB: db},
			args:    args{userName: mocker.userName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			got, err := env.Profile(tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.Profile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			log.Printf("%+v", got)
		})
	}
}
