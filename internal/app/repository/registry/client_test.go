package registry

import (
	"strings"
	"testing"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/faker"
	"github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/pkg"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/brianvoe/gofakeit/v5"
)

func mockNewApp() oauth.App {
	faker.SeedGoFake()

	name := gofakeit.Name()
	slug := strings.Join(strings.Split(name, " "), "-")
	app, err := oauth.NewApp(oauth.BaseApp{
		Name:    name,
		Slug:    slug,
		RepoURL: gofakeit.URL(),
	}, gofakeit.Username())
	if err != nil {
		panic(err)
	}

	return app
}

func TestRandom(t *testing.T) {
	t.Logf("App name: %s", gofakeit.AppName())
	t.Logf("Name: %s", gofakeit.Name())
	t.Logf("Beer name: %s", gofakeit.BeerName())
	t.Logf("Username: %s", gofakeit.Username())
}

func TestEnv_CreateApp(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	app := mockNewApp()

	type args struct {
		app oauth.App
	}
	tests := []struct {
		name    string
		args    args
		want    oauth.App
		wantErr bool
	}{
		{
			name: "create app",
			args: args{
				app: app,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.CreateApp(tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.CreateApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_countApp(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	tests := []struct {
		name    string
		want    int64
		wantErr bool
	}{
		{
			name:    "count app",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.countApp()
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.countApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%d", got)
		})
	}
}

func TestEnv_listApps(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    []oauth.App
		wantErr bool
	}{
		{
			name: "list app",
			args: args{
				p: gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.listApps(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.listApps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_ListApps(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	type args struct {
		p gorest.Pagination
	}
	tests := []struct {
		name    string
		args    args
		want    pkg.PagedList[oauth.App]
		wantErr bool
	}{
		{
			name: "list app",
			args: args{
				p: gorest.NewPagination(1, 10),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.ListApps(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.ListApps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("%v", got)
		})
	}
}

func TestEnv_RetrieveApp(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	app := mockNewApp()

	app, _ = env.CreateApp(app)

	type args struct {
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    oauth.App
		wantErr bool
	}{
		{
			name: "retrieve app",
			args: args{
				clientID: app.ClientID.String(),
			},
			want:    app,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := env.RetrieveApp(tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.RetrieveApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("Env.RetrieveApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_UpdateApp(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	app := mockNewApp()

	app, _ = env.CreateApp(app)

	type args struct {
		app oauth.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update app",
			args: args{
				app: app,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.UpdateApp(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("Env.UpdateApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnv_RemoveApp(t *testing.T) {
	env := NewEnv(db.MockGormSQL())

	app := mockNewApp()

	app, _ = env.CreateApp(app)

	app = app.Remove()

	type args struct {
		app oauth.App
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "remove app",
			args: args{
				app: app,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := env.RemoveApp(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("Env.RemoveApp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
