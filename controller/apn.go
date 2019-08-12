package controller

import (
	"database/sql"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/models/push"
	apn2 "gitlab.com/ftchinese/backyard-api/repository/apn"
	"net/http"
)

type APNRouter struct {
	model apn2.APNEnv
}

func NewAPNRouter(db *sql.DB) APNRouter {
	return APNRouter{
		model: apn2.APNEnv{DB: db},
	}
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
	var d push.TestDevice

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
