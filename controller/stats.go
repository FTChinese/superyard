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
//
// If both `start` and `end` are missing from query parameters, the time range defaults to the past 7 days.
//
// If `start` is missing, it defaults to 7 days earlier before `end`.
// If `end` is missing, it defaults to 7 days later after `start`.
// UTC+08:00 is used rather than UTC time.
//
// - 200 OK with body:
// 	[{
// 		"count": 123,
// 		"date": ""
// 	}]
func (r StatsRouter) DailySignup(w http.ResponseWriter, req *http.Request) {
	start := getQueryParam(req, "start").toString()
	end := getQueryParam(req, "end").toString()

	start, end, err := normalizeTimeRange(start, end)

	if err != nil {
		view.Render(w, util.NewBadRequest("Time format must be YYYY-MM-DD"))

		return
	}

	singups, err := r.model.DailyNewUser(start, end)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(singups))
}
