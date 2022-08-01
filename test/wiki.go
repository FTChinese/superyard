//go:build !production

package test

import (
	"github.com/FTChinese/go-rest/chrono"
	wiki2 "github.com/FTChinese/superyard/internal/pkg/wiki"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"time"
)

func NewArticle() wiki2.Article {
	gofakeit.Seed(time.Now().UnixNano())

	return wiki2.Article{
		ArticleMeta: wiki2.ArticleMeta{
			ID:         0,
			Author:     "weiguo.ni",
			CreatedUTC: chrono.TimeNow(),
			UpdatedUTC: chrono.TimeNow(),
		},
		ArticleInput: wiki2.ArticleInput{
			Title:   gofakeit.Sentence(10),
			Summary: null.StringFrom(gofakeit.Sentence(30)),
			Keyword: null.StringFrom(gofakeit.Word()),
			Body:    null.StringFrom(gofakeit.LoremIpsumParagraph(5, 5, 10, ".")),
		},
	}
}
