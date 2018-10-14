package ftcapi

import (
	"testing"

	"github.com/gosimple/slug"
	"github.com/icrowley/fake"
)

func TestNewApp(t *testing.T) {
	appName := fake.FullName()
	slugName := slug.Make(appName)
	err := devEnv.NewApp(App{
		Name:        appName,
		Slug:        slugName,
		RepoURL:     fake.DomainName(),
		Description: fake.Sentence(),
		OwnedBy:     mockApp.OwnedBy,
	})

	if err != nil {
		t.Error(err)
	}
}

func TestAppRoster(t *testing.T) {
	apps, err := devEnv.AppRoster(1, 20)

	if err != nil {
		t.Error(err)
	}

	t.Log(apps)
}

func TestRetrieveApp(t *testing.T) {
	app, err := devEnv.RetrieveApp(mockApp.Slug)

	if err != nil {
		t.Error(err)
	}

	t.Log(app)
}

func TestUpdateApp(t *testing.T) {
	err := devEnv.UpdateApp(mockApp.Slug, mockApp)

	if err != nil {
		t.Error(err)
	}
}
