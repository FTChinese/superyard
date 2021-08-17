package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/faker"
	oauth2 "github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_CreateApp(t *testing.T) {

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		app oauth2.App
	}

	env := NewEnv(db.MustNewMyDBs(false))

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Create a new app",
			args:    args{app: test.FixedStaff.MustNewOAuthApp()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.CreateApp(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("CreateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_ListApps(t *testing.T) {

	test.NewRepo().MustCreateOAuthApp(test.FixedStaff.MustNewOAuthApp())

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "List Apps",
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

			t.Logf("%s", faker.MustMarshalIndent(got))
		})
	}
}

func TestEnv_RetrieveApp(t *testing.T) {
	app := test.FixedStaff.MustNewOAuthApp()

	test.NewRepo().MustCreateOAuthApp(app)

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Retrieve an app",
			args:    args{clientID: app.ClientID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		app oauth2.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Update an app",
			args:    args{app: app},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := env.UpdateApp(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("UpdateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
