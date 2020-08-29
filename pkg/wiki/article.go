package wiki

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
	"strings"
)

type ArticleInput struct {
	Title   string      `json:"title" db:"title"`
	Summary null.String `json:"summary" db:"summary"`
	Keyword null.String `json:"keyword" db:"keyword"`
	Body    null.String `json:"body" db:"body"`
}

func (i *ArticleInput) Validate() *render.ValidationError {
	i.Title = strings.TrimSpace(i.Title)
	i.Summary.String = strings.TrimSpace(i.Summary.String)
	i.Keyword.String = strings.TrimSpace(i.Keyword.String)
	i.Body.String = strings.TrimSpace(i.Body.String)

	return validator.New("title").Required().Validate(i.Title)
}

// Update creates an Article instance with an existing id.
func (i ArticleInput) Update(id int64) Article {
	return Article{
		ArticleMeta: ArticleMeta{
			ID:         id,
			UpdatedUTC: chrono.TimeNow(),
		},
		ArticleInput: i,
	}
}

// ArticleMeta contains metadata of an article.
type ArticleMeta struct {
	ID         int64       `json:"id" db:"id"`
	Author     string      `json:"author" db:"author"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

// Article contains the full data of an article.
type Article struct {
	ArticleMeta
	ArticleInput
}

// NewArticle creates a new Article based on user input.
func NewArticle(input ArticleInput, author string) Article {
	return Article{
		ArticleMeta: ArticleMeta{
			ID:         0,
			Author:     author,
			CreatedUTC: chrono.TimeNow(),
			UpdatedUTC: chrono.TimeNow(),
		},
		ArticleInput: input,
	}
}

// Update an existing article. Since the request body
// does not contain the article's id, you have to get it
// from the path parameter.
func (a Article) Update(input ArticleInput) Article {
	a.ArticleInput = input
	a.UpdatedUTC = chrono.TimeNow()

	return a
}
