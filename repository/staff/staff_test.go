package staff

import (
	"github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/employee"
	"gitlab.com/ftchinese/backyard-api/test"
	"testing"
)

func createStaff(a employee.Account) employee.Account {
	env := Env{DB: test.DBX}

	if err := env.Create(a); err != nil {
		panic(err)
	}

	return a
}

func TestEnv_Create(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		a employee.Account
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Create a Staff",
			fields:  fields{DB: test.DBX},
			args:    args{a: test.NewStaff().Account()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			if err := env.Create(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_Load(t *testing.T) {

	a := createStaff(test.NewStaff().Account())

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		col   Column
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Load Staff Profile by ID",
			fields: fields{DB: test.DBX},
			args: args{
				col:   ColumnStaffId,
				value: a.ID,
			},
			wantErr: false,
		},
		{
			name:   "Load Staff Profile by Name",
			fields: fields{DB: test.DBX},
			args: args{
				col:   ColumnUserName,
				value: a.UserName,
			},
			wantErr: false,
		},
		{
			name:   "Load Staff Profile by Email",
			fields: fields{DB: test.DBX},
			args: args{
				col:   ColumnEmail,
				value: a.Email,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.Load(tt.args.col, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Staff profile: %+v", got)
		})
	}
}

func TestEnv_List(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "List Staff",
			fields: fields{DB: test.DBX},
			args: args{
				p: gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.List(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("All staff %+v", got)
		})
	}
}
