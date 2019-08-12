package customer

import (
	"database/sql"
	"gitlab.com/ftchinese/backyard-api/model"
	"testing"
)

func TestUserEnv_loadAccount(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		col model.sqlUserCol
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
			fields:  fields{DB: model.db},
			args:    args{col: model.sqlUserColID, val: model.myFtcID},
			wantErr: false,
		},
		{
			name:    "Wechat Account",
			fields:  fields{DB: model.db},
			args:    args{col: model.sqlUserColUnionID, val: model.myUnionID},
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
