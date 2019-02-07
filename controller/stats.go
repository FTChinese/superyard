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

// DailySignUp show how many new users signed up at ftchinese.com everyday.
//
//	GET /stats/signup/daily?start=YYYY-MM-DD&end=YYYY-MM-DD
func (r StatsRouter) DailySignUp(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	start, _ := GetQueryParam(req, "start").ToString()
	end, _ := GetQueryParam(req,"end").ToString()

	log.WithField("trace", "DailySignUp").Infof("Original start and end: %s - %s", start, end)

	period, err := stats.NewPeriod(start, end)
	if err != nil {
		view.Render(w, view.NewBadRequest("Time format must be YYYY-MM-DD"))
		return
	}

	signups, err := r.model.DailyNewUser(period)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(signups))
}
