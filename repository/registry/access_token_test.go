package registry

import (
	"github.com/FTChinese/go-rest"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/oauth"
	"gitlab.com/ftchinese/backyard-api/test"
	"testing"
)

func TestEnv_CreateToken(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}

	env := Env{DB: test.DBX}

	personalToken, err := oauth.NewAccess(oauth.InputKey{
		Description: null.String{},
		CreatedBy:   "weiguo.ni",
		ClientID:    null.String{},
	})

	t.Logf("Personal Access Token: %v", personalToken)

	if err != nil {
		t.Error(err)
	}

	type args struct {
		acc oauth.Access
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Personal Access Token",
			args: args{
				acc: personalToken,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.CreateToken(tt.args.acc)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Token id: %d", got)
		})
	}
}

func TestEnv_ListKeys(t *testing.T) {
	env := Env{DB: test.DBX}

	type args struct {
		by oauth.KeySelector
		p  gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "List Personal Keys",
			args: args{
				by: oauth.KeySelector{
					StaffName: null.StringFrom("weiguo.ni"),
				},
				p: gorest.NewPagination(1, 20),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListKeys(tt.args.by, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%+v", got)
		})
	}
}
