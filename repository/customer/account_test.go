package customer

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/test"
	"testing"
)

func TestEnv_LoadAccountFtc(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		ftcID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve FTC Account",
			fields:  fields{DB: test.DBX},
			args:    args{ftcID: "4f3b1973-f7ee-42b2-a123-3721dadc542f"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.LoadAccountFtc(tt.args.ftcID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAccountFtc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Ftc account: %+v", got)
		})
	}
}

func TestEnv_LoadAccountWx(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		unionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve Wechat Account",
			fields:  fields{DB: test.DBX},
			args:    args{unionID: "a2xlwxV2bpPZIRWXjADTW0XGVkqg"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.LoadAccountWx(tt.args.unionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAccountWx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Wechat Account: %+v", got)
		})
	}
}
