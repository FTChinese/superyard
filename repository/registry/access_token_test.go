package registry

import (
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/superyard/models/oauth"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_CreateToken(t *testing.T) {
	type fields struct {
		DB *sqlx.DB
	}

	env := Env{DB: test.DBX}

	personalToken, err := oauth.NewAccess(oauth.BaseAccess{
		Description: null.String{},
		ClientID:    null.String{},
	}, "weiguo.ni")

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
