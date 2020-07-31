package registry

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ftchinese/superyard/pkg/oauth"
	"gitlab.com/ftchinese/superyard/test"
	"testing"
)

func TestEnv_CreateToken(t *testing.T) {
	s := test.NewStaff()
	app := s.MustNewOAuthApp()

	repo := test.NewRepo()
	repo.MustCreateStaff(s.SignUp())
	repo.MustCreateOAuthApp(app)

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		acc oauth.Access
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Create personal key",
			fields: fields{DB: test.DBX},
			args: args{
				acc: s.MustNewPersonalKey(),
			},
			wantErr: false,
		},
		{
			name: "Create app token",
			fields: fields{
				DB: test.DBX,
			},
			args: args{
				acc: s.MustNewAppToken(app),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
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

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		clientID string
		p        gorest.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "List app access token",
			fields: fields{DB: test.DBX},
			args: args{
				clientID: app.ClientID,
				p:        gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.ListAccessTokens(tt.args.clientID, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAccessTokens() error = %v, wantErr %v", err, tt.wantErr)
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

	type fields struct {
		DB *sqlx.DB
	}
	type args struct {
		owner string
		p     gorest.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Personal keys",
			fields: fields{DB: test.DBX},
			args: args{
				owner: s.UserName,
				p:     gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				DB: tt.fields.DB,
			}
			got, err := env.ListPersonalKeys(tt.args.owner, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListPersonalKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, 1, len(got))

			t.Logf("Personal keys: %+v", got)
		})
	}
}
