package wiki

import (
	"strings"

	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/guregu/null"
)

type ArticleInput struct {
	Title   string      `json:"title" db:"title" gorm:"column:title"`
	Summary null.String `json:"summary" db:"summary" gorm:"column:summary"`
	Keyword null.String `json:"keyword" db:"keyword" gorm:"column:keyword"`
	Body    null.String `json:"body" db:"body" gorm:"column:body"`
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
	ID         int64       `json:"id" db:"id" gorm:"column:id"`
	Author     string      `json:"author" db:"author" gorm:"column:author"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc" gorm:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc" gorm:"updated_utc"`
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

func (a Article) TableName() string {
	return "file_store.wiki"
}
