//go:build !production

package wiki

import (
	"time"

	"github.com/FTChinese/go-rest/chrono"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
)

func MockArticle() Article {
	gofakeit.Seed(time.Now().UnixNano())

	return Article{
		ArticleMeta: ArticleMeta{
			ID:         0,
			Author:     "weiguo.ni",
			CreatedUTC: chrono.TimeNow(),
			UpdatedUTC: chrono.TimeNow(),
		},
		ArticleInput: ArticleInput{
			Title:   gofakeit.Sentence(10),
			Summary: null.StringFrom(gofakeit.Sentence(30)),
			Keyword: null.StringFrom(gofakeit.Word()),
			Body:    null.StringFrom(gofakeit.LoremIpsumParagraph(5, 5, 10, ".")),
		},
	}
}
