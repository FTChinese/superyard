package controller

import (
	"database/sql"
	"net/http"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/subscription"
	"gitlab.com/ftchinese/backyard-api/util"
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
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	page, err := getQueryParam(req, "page").toInt()

	if err != nil {
		page = 1
	}

	promos, err := sr.model.ListPromo(page, 5)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(promos))
}

// CreateSchedule saves the schedule part of a promotion compaign.
//
//	POST /subscripiton/promos
//
// Request body is type subscription.Schedule without `id` field.
func (sr SubsRouter) CreateSchedule(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var sch subscription.Schedule
	if err := util.Parse(req.Body, &sch); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	sch.Sanitize()

	if r := sch.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	id, err := sr.model.NewSchedule(sch, userName)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewResponse().SetBody(map[string]int64{
		"id": id,
	}))
}

// GetPromo loads a piece of promotion.
//
// GET /subscription/promos/{id}
func (sr SubsRouter) GetPromo(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	promo, err := sr.model.RetrievePromo(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewResponse().SetBody(promo))
}

// RemovePromo deletes a record.
//
// DELETE `/subscription/promos/{id}`
func (sr SubsRouter) RemovePromo(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	err = sr.model.DisablePromo(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewNoContent())
}

// SetPromoPricing saves/updates a promotion's pricing plans.
//
// PATCH /subscription/promos/{id}/pricing
func (sr SubsRouter) SetPromoPricing(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	var plans subscription.Pricing

	if err := util.Parse(req.Body, &plans); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	err = sr.model.SavePricing(id, plans)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewNoContent())
}

// SetPromoBanner saves/updates a promotion's banner content
//
// POST /subscription/promos/{id}/banner
func (sr SubsRouter) SetPromoBanner(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	var banner subscription.Banner
	if err := util.Parse(req.Body, &banner); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	banner.Sanitize()

	if r := banner.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))

		return
	}

	err = sr.model.SaveBanner(id, banner)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewNoContent())
}
