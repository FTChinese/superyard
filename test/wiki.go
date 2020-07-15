package test

import (
	"github.com/brianvoe/gofakeit/v5"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/superyard/pkg/wiki"
	"time"
)

func NewArticle() wiki.Article {
	gofakeit.Seed(time.Now().UnixNano())

	return wiki.Article{
		ArticleTeaser: wiki.ArticleTeaser{
			ArticleMeta: wiki.NewArticleMeta("weiguo.ni"),
			Title:       gofakeit.Sentence(10),
			Summary:     null.StringFrom(gofakeit.Sentence(30)),
			Keyword:     null.StringFrom(gofakeit.Word()),
		},
		Body: gofakeit.LoremIpsumParagraph(5, 5, 10, "."),
	}
}
