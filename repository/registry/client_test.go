package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ftchinese/superyard/pkg/oauth"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_CreateApp(t *testing.T) {

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		app oauth.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Create a new app",
			fields:  fields{DB: test.DBX},
			args:    args{app: test.FixedStaff.MustNewOAuthApp()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			if err := env.CreateApp(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("CreateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_ListApps(t *testing.T) {

	test.NewRepo().MustCreateOAuthApp(test.FixedStaff.MustNewOAuthApp())

	env := Env{DB: test.DBX}

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
			name: "List Apps",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				p: gorest.NewPagination(1, 10),
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
			assert.NotZero(t, len(got))
		})
	}
}

func TestEnv_RetrieveApp(t *testing.T) {
	app := test.FixedStaff.MustNewOAuthApp()

	test.NewRepo().MustCreateOAuthApp(app)

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		clientID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve an app",
			fields:  fields{DB: test.DBX},
			args:    args{clientID: app.ClientID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.RetrieveApp(tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetrieveApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.ClientID, app.ClientID)
		})
	}
}

func TestEnv_UpdateApp(t *testing.T) {
	app := test.FixedStaff.MustNewOAuthApp()

	test.NewRepo().MustCreateOAuthApp(app)

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		app oauth.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Update an app",
			fields:  fields{DB: test.DBX},
			args:    args{app: app},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			if err := env.UpdateApp(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("UpdateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
