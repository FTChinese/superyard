//go:build !production

package test

import (
	oauth2 "github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"time"
)

func genOAuthApp() oauth2.BaseApp {
	gofakeit.Seed(time.Now().UnixNano())

	return oauth2.BaseApp{
		Name:        gofakeit.Name(),
		Slug:        gofakeit.Username(),
		RepoURL:     gofakeit.URL(),
		Description: null.StringFrom(gofakeit.Sentence(20)),
		HomeURL:     null.StringFrom(gofakeit.URL()),
		CallbackURL: null.StringFrom(gofakeit.URL()),
	}
}
