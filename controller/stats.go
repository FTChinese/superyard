package controller

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/stats"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// StatsRouter responds to requests for statistic data.
type StatsRouter struct {
	model stats.Env
}

// NewStatsRouter creates a new instance of StatsRouter
func NewStatsRouter(db *sql.DB) StatsRouter {
	model := stats.Env{DB: db}

	return StatsRouter{
		model: model,
	}
}

// DailySignup show how many new users signed up at ftchinese.com everyday.
//
//	GET /stats/signup/daily?start=YYYY-MM-DD&end=YYYY-MM-DD
func (r StatsRouter) DailySignup(w http.ResponseWriter, req *http.Request) {
	start := getQueryParam(req, "start").toString()
	end := getQueryParam(req, "end").toString()

	start, end, err := normalizeTimeRange(start, end)

	if err != nil {
		view.Render(w, util.NewBadRequest("Time format must be YYYY-MM-DD"))

		return
	}

	signups, err := r.model.DailyNewUser(start, end)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(signups))
}
