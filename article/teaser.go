package article

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

type Teaser struct {
	ArticleID  string      `json:"id"`
	Title      string      `json:"title"`
	CoverURL   null.String `json:"coverUrl"`
	Standfirst string      `json:"standfirst"`
	Author     string      `json:"author"`
	Tags       []string    `json:"tags"`
	CreatedAt  chrono.Time `json:"createdAt"`
	UpdateAt   chrono.Time `json:"updatedAt"`
}
