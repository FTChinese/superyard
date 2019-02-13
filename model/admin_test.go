package model

import (
	"database/sql"
	"log"
	"testing"

	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
)

func TestAdminEnv_CreateAccount(t *testing.T) {
	mock := newMockStaff()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		a        staff.Account
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Create New Account for Staff",
			fields:  fields{DB: db},
			args:    args{a: mock.account()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AdminEnv{
				DB: tt.fields.DB,
			}
			if err := env.CreateAccount(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("AdminEnv.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminEnv_ListAccounts(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		p util.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "List Accounts",
			fields:  fields{DB: db},
			args:    args{p: util.NewPagination(1, 20)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AdminEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListAccounts(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminEnv.ListAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			log.Printf("%+v", got)
		})
	}
}

func TestAdminEnv_UpdateAccount(t *testing.T) {
	mock := newMockStaff()
	mock.createAccount()
	log.Printf("Original account: %+v", mock.account())

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName string
		a        staff.Account
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Update Account",
			fields: fields{DB: db},
			args: args{
				userName: mock.userName,
				a:        newMockStaff().account(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AdminEnv{
				DB: tt.fields.DB,
			}
			if err := env.UpdateAccount(tt.args.userName, tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("AdminEnv.UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminEnv_RemoveStaff(t *testing.T) {
	mock := newMockStaff()
	mock.createAccount()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName  string
		revokeVIP bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Deactivate Account",
			fields: fields{DB: db},
			args: args{
				userName:  mock.userName,
				revokeVIP: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AdminEnv{
				DB: tt.fields.DB,
			}
			if err := env.RemoveStaff(tt.args.userName, tt.args.revokeVIP); (err != nil) != tt.wantErr {
				t.Errorf("AdminEnv.RemoveStaff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminEnv_ActivateStaff(t *testing.T) {
	mock := newMockStaff()
	mock.createAccount()

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
			name:    "Activate Staff",
			fields:  fields{DB: db},
			args:    args{userName: mock.userName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AdminEnv{
				DB: tt.fields.DB,
			}
			if err := env.ActivateStaff(tt.args.userName); (err != nil) != tt.wantErr {
				t.Errorf("AdminEnv.ActivateStaff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminEnv_ListVIP(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "List VIP",
			fields:  fields{DB: db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AdminEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListVIP()
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminEnv.ListVIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			log.Printf("%+v", got)
		})
	}
}

func TestAdminEnv_GrantVIP(t *testing.T) {
	mUser := newMockUser()
	u := mUser.createUser()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		myftID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Grant VIP",
			fields:  fields{DB: db},
			args:    args{myftID: u.UserID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AdminEnv{
				DB: tt.fields.DB,
			}
			if err := env.GrantVIP(tt.args.myftID); (err != nil) != tt.wantErr {
				t.Errorf("AdminEnv.GrantVIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminEnv_RevokeVIP(t *testing.T) {
	mUser := newMockUser()
	u := mUser.createUser()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		myftID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Revoke VIP",
			fields:  fields{DB: db},
			args:    args{myftID: u.UserID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := AdminEnv{
				DB: tt.fields.DB,
			}
			if err := env.RevokeVIP(tt.args.myftID); (err != nil) != tt.wantErr {
				t.Errorf("AdminEnv.RevokeVIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
