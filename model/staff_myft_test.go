package model

import (
	"database/sql"
	"log"
	"testing"

	"gitlab.com/ftchinese/backyard-api/user"
)

func TestStaffEnv_AddMyft(t *testing.T) {
	uMock := newMockUser()
	uMock.createUser()

	sMock := newMockStaff()
	sMock.createAccount()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		staffName string
		l         user.Login
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Add Myft",
			fields: fields{DB: db},
			args: args{
				staffName: sMock.userName,
				l:         uMock.login(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			if err := env.AddMyft(tt.args.staffName, tt.args.l); (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.AddMyft() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffEnv_ListMyft(t *testing.T) {
	sMock := newMockStaff()

	for i := 0; i < 5; i++ {
		u := newMockUser().createUser()
		sMock.createMyft(u)
	}

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		staffName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "List Myft",
			fields:  fields{DB: db},
			args:    args{staffName: sMock.userName},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListMyft(tt.args.staffName)
			if (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.ListMyft() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}

func TestStaffEnv_DeleteMyft(t *testing.T) {
	u := newMockUser().createUser()

	sMock := newMockStaff()
	myft := sMock.createMyft(u)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		staffName string
		myftID    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Delete Myft",
			fields:  fields{DB: db},
			args:    args{staffName: sMock.userName, myftID: myft.MyftID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := StaffEnv{
				DB: tt.fields.DB,
			}
			if err := env.DeleteMyft(tt.args.staffName, tt.args.myftID); (err != nil) != tt.wantErr {
				t.Errorf("StaffEnv.DeleteMyft() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
