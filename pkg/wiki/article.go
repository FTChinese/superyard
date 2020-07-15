package wiki

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

type BaseArticle struct {
	ID         int64       `json:"id" db:"id"`
	Title      string      `json:"title" db:"title"`
	Author     string      `json:"author" db:"author"`
	Summary    null.String `json:"summary" db:"summary"`
	Keyword    null.String `json:"keyword" db:"keyword"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	UpdatedUTC chrono.Time `json:"updatedUtc" db:"updated_utc"`
}

type Article struct {
	BaseArticle
	Body string `json:"body" db:"body"`
}
