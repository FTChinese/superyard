package controller

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/stats"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// StatsRouter handles request for statistics
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

// DailySignup outputs new user for everyday
func (r StatsRouter) DailySignup(w http.ResponseWriter, req *http.Request) {
	start := getQueryParam(req, "start").toString()
	end := getQueryParam(req, "end").toString()

	if start == "" {
		start = util.DateNow(8 * 60 * 60)
	}

	if end == "" {
		end = util.DateNow(8*60*60 + 7*24*60*60)
	}

	singups, err := r.model.DailyNewUser(start, end)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(singups))
}
