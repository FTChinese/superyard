package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/model"
	"net/http"
)

type APNRouter struct {
	model model.APNEnv
}

func NewAPNRouter(db *sql.DB) APNRouter {
	return APNRouter{
		model: model.APNEnv{DB: db},
	}
}

func (router APNRouter) LatestStoryList(w http.ResponseWriter, req *http.Request) {
	teasers, err := router.model.LatestStoryList()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(teasers))
}

func (router APNRouter) StoryTeaser(w http.ResponseWriter, req *http.Request) {
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

func (router APNRouter) VideoTeaser(w http.ResponseWriter, req *http.Request) {
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

func (router APNRouter) GalleryTeaser(w http.ResponseWriter, req *http.Request) {
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

func (router APNRouter) InteractiveTeaser(w http.ResponseWriter, req *http.Request) {
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
