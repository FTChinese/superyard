package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/model"
	"net/http"
)

type ContentRouter struct {
	model model.ArticleEnv
}

func NewContentRouter(db *sql.DB) ContentRouter {
	return ContentRouter{
		model: model.ArticleEnv{DB: db},
	}
}

func (router ContentRouter) LatestStoryList(w http.ResponseWriter, req *http.Request) {
	teasers, err := router.model.LatestStoryList()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(teasers))
}

func (router ContentRouter) StoryTeaser(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	teaser, err := router.model.FindStory(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(teaser))
}

func (router ContentRouter) VideoTeaser(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	teaser, err := router.model.FindVideo(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(teaser))
}

func (router ContentRouter) GalleryTeaser(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	teaser, err := router.model.FindGallery(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(teaser))
}

func (router ContentRouter) InteractiveTeaser(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	teaser, err := router.model.FindInteractive(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(teaser))
}
