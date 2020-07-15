package wiki

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

// ArticleMeta contains metadata of an article.
type ArticleMeta struct {
	ID         int64       `json:"id" db:"id"`
	Author     string      `json:"author" db:"author"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

func NewArticleMeta(author string) ArticleMeta {
	return ArticleMeta{
		ID:         0,
		Author:     author,
		CreatedUTC: chrono.TimeNow(),
		UpdatedUTC: chrono.TimeNow(),
	}
}

// ArticleOverview is used as an item of a list of articles.
type ArticleTeaser struct {
	ArticleMeta
	Title   string      `json:"title" db:"title"`
	Summary null.String `json:"summary" db:"summary"`
	Keyword null.String `json:"keyword" db:"keyword"`
}

// Article contains the full data of an article.
type Article struct {
	ArticleTeaser
	Body string `json:"body" db:"body"`
}

// NewArticle creates a new Article based on user input.
func NewArticle(input Article, author string) Article {
	input.ArticleMeta = NewArticleMeta(author)

	return input
}

// Update an existing article. Since the request body
// does not contain the article's id, you have to get it
// from the path parameter.
func (a Article) Update(id int64) Article {
	a.ID = id
	a.UpdatedUTC = chrono.TimeNow()

	return a
}
