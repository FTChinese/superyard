package admin

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_List(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
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
			name:   "ListStaff Staff",
			fields: fields{DB: test.DBX},
			args: args{
				p: util.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.ListStaff(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("All staff %+v", got)
		})
	}
}
