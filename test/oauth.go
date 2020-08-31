package test

import (
	"github.com/FTChinese/superyard/pkg/oauth"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"time"
)

func genOAuthApp() oauth.BaseApp {
	gofakeit.Seed(time.Now().UnixNano())

	return oauth.BaseApp{
		Name:        gofakeit.Name(),
		Slug:        gofakeit.Username(),
		RepoURL:     gofakeit.URL(),
		Description: null.StringFrom(gofakeit.Sentence(20)),
		HomeURL:     null.StringFrom(gofakeit.URL()),
		CallbackURL: null.StringFrom(gofakeit.URL()),
	}
}
