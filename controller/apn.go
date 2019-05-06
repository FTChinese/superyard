package controller

import (
	"database/sql"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/apn"
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

func (router APNRouter) ListMessages(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request if query string cannot be parsed.
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	msgs, err := router.model.ListMessage(pagination)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(msgs))
}

func (router APNRouter) LoadTimezones(w http.ResponseWriter, req *http.Request) {
	tz, err := router.model.TimeZoneDist()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(tz))
}

func (router APNRouter) LoadDeviceDist(w http.ResponseWriter, req *http.Request) {
	d, err := router.model.DeviceDist()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(d))
}

func (router APNRouter) LoadInvalidDist(w http.ResponseWriter, req *http.Request) {
	d, err := router.model.InvalidDist()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(d))
}

func (router APNRouter) CreateTestDevice(w http.ResponseWriter, req *http.Request) {
	var d apn.TestDevice

	if err := gorest.ParseJSON(req.Body, &d); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	err := router.model.CreateTestDevice(d)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

func (router APNRouter) ListTestDevice(w http.ResponseWriter, req *http.Request) {
	d, err := router.model.ListTestDevice()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(d))
}

func (router APNRouter) RemoveTestDevice(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	err = router.model.RemoveTestDevice(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}
