package controller

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/repository/paywall"
	"net/http"

	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/models/promo"
)

// PromoRouter handles request for subs related data.
type PromoRouter struct {
	model paywall.PromoEnv
}

// NewPromoRouter creates a new instance of SubscriptionRouter
func NewPromoRouter(db *sqlx.DB) PromoRouter {
	return PromoRouter{
		model: paywall.PromoEnv{DB: db},
	}
}

// CreateSchedule saves the schedule part of a promotion campaign.
//
//	POST /subs/schedule
//
// Input {id: number, name: string, description: null | string, startAt: string, endAt: string}
func (router PromoRouter) CreateSchedule(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var sch promo.Schedule
	if err := gorest.ParseJSON(req.Body, &sch); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	sch.Sanitize()

	if r := sch.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	id, err := router.model.NewSchedule(sch, userName)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewResponse().SetBody(map[string]int64{
		"id": id,
	}))
}

// SetPricingPlans saves/updates a promotion's pricing plans.
//
// PATCH /subs/schedule/{id}/pricing
func (router PromoRouter) SetPricingPlans(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	var plans promo.Pricing

	if err := gorest.ParseJSON(req.Body, &plans); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	err = router.model.SavePlans(id, plans)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// SetPromoBanner saves/updates a promotion's banner content
//
// POST /subs/schedule/{id}/banner
func (router PromoRouter) SetBanner(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var banner promo.Banner
	if err := gorest.ParseJSON(req.Body, &banner); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	banner.Sanitize()

	if r := banner.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))

		return
	}

	err = router.model.SaveBanner(id, banner)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// ListPromos list promotion schedules by page.
//
// GET /subs/promos?page=<int>&per_page=<number>
func (router PromoRouter) ListPromos(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	pagination := gorest.GetPagination(req)

	promos, err := router.model.ListPromos(pagination)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(promos))
}

// GetPromo loads a piece of promotion.
//
// GET /subs/promos/{id}
func (router PromoRouter) LoadPromo(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	promo, err := router.model.LoadPromo(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(promo))
}

// RemovePromo deletes a record.
//
// DELETE `/subs/promos/{id}`
func (router PromoRouter) DisablePromo(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	err = router.model.DisablePromo(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}
