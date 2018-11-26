package controller

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/subscription"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// SubsRouter handles request for subscription related data.
type SubsRouter struct {
	model subscription.Env
}

// NewSubsRouter creates a new isntance of SubscriptionRouter
func NewSubsRouter(db *sql.DB) SubsRouter {
	model := subscription.Env{DB: db}

	return SubsRouter{
		model: model,
	}
}

// ListPromos list promotion schedules by page.
//
// GET `/subscription/promos?page=<int>`
func (sr SubsRouter) ListPromos(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	page, err := getQueryParam(req, "page").toInt()

	if err != nil {
		page = 1
	}

	promos, err := sr.model.ListPromo(page, 10)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(promos))
}

// CreateSchedule saves a new schedule.
//
// POST `/subscripiton/proms`
func (sr SubsRouter) CreateSchedule(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var sch subscription.Schedule
	if err := util.Parse(req.Body, &sch); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	sch.Sanitize()

	if r := sch.Validate(); r != nil {
		view.Render(w, util.NewUnprocessable(r))
		return
	}

	sch.CreatedBy = userName
	id, err := sr.model.NewSchedule(sch)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewResponse().SetBody(map[string]int64{
		"id": id,
	}))
}

// GetPromo loads a piece of promotion.
//
// GET /subscription/promos/{id}
func (sr SubsRouter) GetPromo(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	promo, err := sr.model.RetrievePromo(id)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewResponse().SetBody(promo))
}

// RemovePromo deletes a record.
//
// DELETE `/subscription/promos/{id}`
func (sr SubsRouter) RemovePromo(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	err = sr.model.DeletePromo(id)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewNoContent())
}

// SetPromoPricing saves a promotion's pricing plans.
//
// PATCH /subscription/promos/{id}/pricing
func (sr SubsRouter) SetPromoPricing(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	var plans map[string]subscription.Plan

	if err := util.Parse(req.Body, &plans); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	err = sr.model.SavePricing(id, plans)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewNoContent())
}

// SetPromoBanner saves a promotion's banner content
//
// POST /subscription/promos/{id}/banner
func (sr SubsRouter) SetPromoBanner(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	var banner subscription.Banner
	if err := util.Parse(req.Body, &banner); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	banner.Sanitize()

	if r := banner.Validate(); r != nil {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err = sr.model.SaveBanner(id, banner)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewNoContent())
}
