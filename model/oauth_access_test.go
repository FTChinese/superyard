package model

import (
	"database/sql"
	"log"
	"reflect"
	"testing"

	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/oauth"
	"gitlab.com/ftchinese/backyard-api/util"
)

func TestOAuthEnv_SaveAppAccess(t *testing.T) {
	m := newMockStaff()

	app := m.createApp()

	token, _ := oauth.NewToken()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		token    string
		clientID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Save App Access",
			fields:  fields{DB: db},
			args:    args{token: token, clientID: app.ClientID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			got, err := env.SaveAppAccess(tt.args.token, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.SaveAppAccess() error = %v, wantErr %v", err, tt.wantErr)
			}

			log.Print(got)
		})
	}
}

func TestOAuthEnv_ListAppAccess(t *testing.T) {
	m := newMockStaff()
	app := m.createApp()

	m.createAppAccess(app)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		slug string
		p    util.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "List App Access",
			fields: fields{DB: db},
			args: args{
				slug: app.Slug,
				p:    util.NewPagination(1, 20)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListAppAccess(tt.args.slug, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.ListAppAccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Print(got)
		})
	}
}

func TestOAuthEnv_RemoveAppAccess(t *testing.T) {
	m := newMockStaff()
	app := m.createApp()
	acc := m.createAppAccess(app)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		clientID string
		id       int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "List App Access",
			fields: fields{DB: db},
			args: args{
				clientID: app.ClientID,
				id:       acc.ID},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			if err := env.RemoveAppAccess(tt.args.clientID, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.RemoveAppAccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOAuthEnv_FindMyftID(t *testing.T) {
	m := newMockUser()
	u := m.createUser()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		email null.String
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   null.String
	}{
		{
			name:   "Find User ID",
			fields: fields{DB: db},
			args:   args{email: null.StringFrom(u.Email)},
			want:   null.StringFrom(m.userID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			if got := env.FindMyftID(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OAuthEnv.FindMyftID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOAuthEnv_SavePersonalToken(t *testing.T) {
	mUser := newMockUser()
	mUser.createUser()

	mStaff := newMockStaff()
	mStaff.createAccount()

	acc := mStaff.personalAccess()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		acc    oauth.PersonalAccess
		myftID null.String
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Save Personal Token",
			fields: fields{DB: db},
			args: args{
				acc:    acc,
				myftID: null.StringFrom(mUser.userID),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			got, err := env.SavePersonalToken(tt.args.acc, tt.args.myftID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.SavePersonalToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			log.Print(got)
		})
	}
}

func TestOAuthEnv_ListPersonalTokens(t *testing.T) {
	mUser := newMockUser()
	u := mUser.createUser()

	mStaff := newMockStaff()
	mStaff.createPersonalToken(u)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName string
		p        util.Pagination
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "List Personal Token",
			fields:  fields{DB: db},
			args:    args{userName: mStaff.userName, p: util.NewPagination(1, 20)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListPersonalTokens(tt.args.userName, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.ListPersonalTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			log.Printf("%+v", got)
		})
	}
}

func TestOAuthEnv_RemovePersonalToken(t *testing.T) {
	mUser := newMockUser()
	u := mUser.createUser()

	mStaff := newMockStaff()
	acc := mStaff.createPersonalToken(u)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		userName string
		id       int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Remove Personal Token",
			fields: fields{DB: db},
			args:   args{userName: mStaff.userName, id: acc.ID},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			if err := env.RemovePersonalToken(tt.args.userName, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.RemovePersonalToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
