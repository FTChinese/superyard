package controller

import (
	"database/sql"
	"net/http"

	"github.com/FTChinese/go-rest/view"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/stats"
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

	start := req.FormValue("start")
	end := req.FormValue("end")

	log.WithField("location", "DailySignup").Infof("Original start and end: %s - %s", start, end)

	start, end, err := normalizeTimeRange(start, end)

	log.WithField("location", "DailySignup").Infof("Normalized start and end: %s - %s", start, end)

	if err != nil {
		view.Render(w, view.NewBadRequest("Time format must be YYYY-MM-DD"))

		return
	}

	signups, err := r.model.DailyNewUser(start, end)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(signups))
}
