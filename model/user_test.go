package model

import (
	"database/sql"
	"testing"
)

func TestUserEnv_loadAccount(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		col sqlUserCol
		val string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "FTC Account by ID",
			fields:  fields{DB: db},
			args:    args{col: sqlUserColID, val: myFtcID},
			wantErr: false,
		},
		{
			name:    "Wechat Account",
			fields:  fields{DB: db},
			args:    args{col: sqlUserColUnionID, val: myUnionID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := UserEnv{
				DB: tt.fields.DB,
			}
			got, err := env.loadAccount(tt.args.col, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserEnv.loadAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Account: %+v", got)
		})
	}
}
