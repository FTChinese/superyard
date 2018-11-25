package controller

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/subscription"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// SubscriptionRouter handles request for subscription related data.
type SubscriptionRouter struct {
	model subscription.Env
}

// NewSubsRouter creates a new isntance of SubscriptionRouter
func NewSubsRouter(db *sql.DB) SubscriptionRouter {
	model := subscription.Env{DB: db}

	return SubscriptionRouter{
		model: model,
	}
}

// CreateSchedule saves a new schedule.
//
// POST `/subscripiton/plans/new`
func (sr SubscriptionRouter) CreateSchedule(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var sch subscription.Promotion
	if err := util.Parse(req.Body, &sch); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// TODO: sanitize, validate

	sch.CreatedBy = userName
	err := sr.model.NewPromo(sch)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, "plans"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// ListSchedules send all schedules.
//
// GET `/subscription/plans?page=<int>`
func (sr SubscriptionRouter) ListSchedules(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	page, err := getQueryParam(req, "page").toInt()

	if err != nil {
		page = 1
	}

	schedules, err := sr.model.ListSchedules(page, 10)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(schedules))
}

// RemoveSchedule deletes a record.
//
// DELETE `/subscription/plans/delete/:id`
func (sr SubscriptionRouter) RemoveSchedule(w http.ResponseWriter, req *http.Request) {
	id, err := getURLParam(req, "id").toInt()

	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	err = sr.model.DeleteSchedule(id)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}
