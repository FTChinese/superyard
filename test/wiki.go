package test

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/superyard/pkg/wiki"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"time"
)

func NewArticle() wiki.Article {
	gofakeit.Seed(time.Now().UnixNano())

	return wiki.Article{
		ArticleMeta: wiki.ArticleMeta{
			ID:         0,
			Author:     "weiguo.ni",
			CreatedUTC: chrono.TimeNow(),
			UpdatedUTC: chrono.TimeNow(),
		},
		ArticleInput: wiki.ArticleInput{
			Title:   gofakeit.Sentence(10),
			Summary: null.StringFrom(gofakeit.Sentence(30)),
			Keyword: null.StringFrom(gofakeit.Word()),
			Body:    null.StringFrom(gofakeit.LoremIpsumParagraph(5, 5, 10, ".")),
		},
	}
}
