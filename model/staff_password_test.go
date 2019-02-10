package model

import (
	"database/sql"
	"log"
	"testing"

	"gitlab.com/ftchinese/backyard-api/staff"
)

func TestStaffEnv_IsPasswordMatched(t *testing.T) {
	m := newMockStaff()
	m.createAccount()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Is Password Matched",
			fields:  fields{DB: db},
			args:    args{userName: m.userName, password: m.password},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			got, err := env.IsPasswordMatched(tt.args.userName, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.IsPasswordMatched() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StaffEnv.IsPasswordMatched() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStaffEnv_SavePwResetToken(t *testing.T) {
	m := newMockStaff()
	a := m.createAccount()

	th, _ := a.TokenHolder()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		h staff.TokenHolder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Save Password Reset Token",
			fields:  fields{DB: db},
			args:    args{h: th},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			if err := env.SavePwResetToken(tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.SavePwResetToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffEnv_VerifyResetToken(t *testing.T) {
	m := newMockStaff()
	th := m.createPwToken()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Verify Password Reset Token",
			fields:  fields{DB: db},
			args:    args{token: th.GetToken()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			got, err := env.VerifyResetToken(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.VerifyResetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}

func TestStaffEnv_ResetPassword(t *testing.T) {
	m := newMockStaff()
	th := m.createPwToken()

	pr := staff.PasswordReset{
		Token:    th.GetToken(),
		Password: genPassword(),
	}

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		r staff.PasswordReset
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Reset Password",
			fields: fields{DB: db},
			args:   args{r: pr},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			if err := env.ResetPassword(tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffEnv_UpdatePassword(t *testing.T) {
	m := newMockStaff()
	m.createAccount()

	p := staff.Password{
		Old: m.password,
		New: genPassword(),
	}

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName string
		p        staff.Password
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Update Password",
			fields:  fields{DB: db},
			args:    args{userName: m.userName, p: p},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			if err := env.UpdatePassword(tt.args.userName, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
