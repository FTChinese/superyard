package registry

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/superyard/models/builder"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_ListApps(t *testing.T) {

	env := Env{DB: test.DBX}

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		p builder.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "List Apps",
			args: args{
				p: builder.NewPagination(1, 10),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListApps(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListApps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", got)
		})
	}
}
