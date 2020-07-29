package test

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/android"
	"time"
)

func NewRelease() android.Release {
	gofakeit.Seed(time.Now().UnixNano())

	return android.Release{
		ReleaseInput: android.ReleaseInput{
			VersionName: semanticVersion(),
			VersionCode: Rand.Int63n(1000),
			Body:        null.StringFrom(gofakeit.Sentence(10)),
			ApkURL:      gofakeit.URL(),
		},
		CreatedAt: chrono.TimeNow(),
		UpdatedAt: chrono.TimeNow(),
	}
}
