package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/oauth"
	"github.com/FTChinese/superyard/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_CreateToken(t *testing.T) {
	s := test.NewStaff()
	app := s.MustNewOAuthApp()

	repo := test.NewRepo()
	repo.MustCreateStaff(s.SignUp())
	repo.MustCreateOAuthApp(app)

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		acc oauth.Access
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create personal key",
			args: args{
				acc: s.MustNewPersonalKey(),
			},
			wantErr: false,
		},
		{
			name: "Create app token",
			args: args{
				acc: s.MustNewAppToken(app),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.CreateToken(tt.args.acc)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.NotZero(t, got)
		})
	}
}

func TestEnv_ListAccessTokens(t *testing.T) {
	s := test.NewStaff()
	app := s.MustNewOAuthApp()

	repo := test.NewRepo()
	repo.MustCreateStaff(s.SignUp())
	repo.MustCreateOAuthApp(app)

	repo.MustInsertAccessToken(s.MustNewPersonalKey())
	repo.MustInsertAccessToken(s.MustNewAppToken(app))

	env := NewEnv(db.MustNewMyDBs(false))

	type args struct {
		clientID string
		p        gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "List app access token",
			args: args{
				clientID: app.ClientID,
				p:        gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListAppTokens(tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAppTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, 1, len(got))

			t.Logf("App tokens: %+v", got)
		})
	}
}

func TestEnv_ListPersonalKeys(t *testing.T) {
	s := test.NewStaff()
	app := s.MustNewOAuthApp()

	repo := test.NewRepo()
	repo.MustCreateStaff(s.SignUp())
	repo.MustCreateOAuthApp(app)

	repo.MustInsertAccessToken(s.MustNewPersonalKey())
	repo.MustInsertAccessToken(s.MustNewAppToken(app))

	env := NewEnv(db.MustNewMyDBs(false))

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		owner string
		p     gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Personal keys",
			args: args{
				owner: s.UserName,
				p:     gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListPersonalKeys(tt.args.owner)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPersonalKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, 1, len(got))

			t.Logf("Personal keys: %+v", got)
		})
	}
}
