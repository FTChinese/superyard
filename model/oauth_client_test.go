package model

import (
	"database/sql"
	"log"
	"testing"

	"gitlab.com/ftchinese/backyard-api/oauth"
	"gitlab.com/ftchinese/backyard-api/util"
)

func TestOAuthEnv_SaveApp(t *testing.T) {
	mStaff := newMockStaff()
	mStaff.createAccount()

	app := mStaff.app()

	type fields struct {
		DB *sql.DB
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
			name:    "Save App",
			fields:  fields{DB: db},
			args:    args{app: app},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			if err := env.SaveApp(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.SaveApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOAuthEnv_ListApps(t *testing.T) {
	m := newMockStaff()
	m.createApp()

	type fields struct {
		DB *sql.DB
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
			name:    "List App",
			fields:  fields{DB: db},
			args:    args{p: util.NewPagination(1, 20)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			got, err := env.ListApps(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.ListApps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}

func TestOAuthEnv_LoadApp(t *testing.T) {
	m := newMockStaff()
	app := m.createApp()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		slug string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Load App",
			fields:  fields{DB: db},
			args:    args{slug: app.Slug},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			got, err := env.LoadApp(tt.args.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.LoadApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			log.Printf("%+v", got)
		})
	}
}

func TestOAuthEnv_UpdateApp(t *testing.T) {
	m := newMockStaff()
	app := m.createApp()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		slug string
		app  oauth.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Update App",
			fields:  fields{DB: db},
			args:    args{slug: app.Slug, app: app},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			if err := env.UpdateApp(tt.args.slug, tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.UpdateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOAuthEnv_FindClientID(t *testing.T) {
	m := newMockStaff()
	app := m.createApp()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		slug string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Find Client ID",
			fields:  fields{DB: db},
			args:    args{slug: app.Slug},
			want:    app.ClientID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			got, err := env.FindClientID(tt.args.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.FindClientID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OAuthEnv.FindClientID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOAuthEnv_RemoveApp(t *testing.T) {
	m := newMockStaff()
	app := m.createApp()

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		clientID string
		owner    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Remove App",
			fields:  fields{DB: db},
			args:    args{clientID: app.ClientID, owner: app.OwnedBy},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			if err := env.RemoveApp(tt.args.clientID, tt.args.owner); (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.RemoveApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOAuthEnv_TransferApp(t *testing.T) {
	m := newMockStaff()
	app := m.createApp()

	m2 := newMockStaff()
	m2.createAccount()

	o := oauth.Ownership{
		SlugName: app.Slug,
		NewOwner: m2.userName,
		OldOwner: m.userName,
	}

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		o oauth.Ownership
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Transfer App",
			fields:  fields{DB: db},
			args:    args{o: o},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := OAuthEnv{
				DB: tt.fields.DB,
			}
			if err := env.TransferApp(tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("OAuthEnv.TransferApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
