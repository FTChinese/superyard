package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/model"
	"net/http"
)

type ArticleRouter struct {
	model model.ArticleEnv
}

func NewArticleRouter(db *sql.DB) ArticleRouter {
	return ArticleRouter{
		model: model.ArticleEnv{DB: db},
	}
}

func (router ArticleRouter) LatestStoryList(w http.ResponseWriter, req *http.Request) {
	teasers, err := router.model.LatestCover()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(teasers))
}
