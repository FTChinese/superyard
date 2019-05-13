package controller

import (
	"database/sql"
	gorest "github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/model"
	"gitlab.com/ftchinese/backyard-api/subs"
	"net/http"
	"time"

	"github.com/FTChinese/go-rest/view"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/stats"
)

// StatsRouter responds to requests for statistic data.
type StatsRouter struct {
	model model.StatsEnv
}

// NewStatsRouter creates a new instance of StatsRouter
func NewStatsRouter(db *sql.DB) StatsRouter {

	return StatsRouter{
		model: model.StatsEnv{DB: db},
	}
}

// DailySignUp show how many new users signed up at ftchinese.com everyday.
//
//	GET /stats/signup/daily?start=YYYY-MM-DD&end=YYYY-MM-DD
func (router StatsRouter) DailySignUp(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	start, _ := gorest.GetQueryParam(req, "start").ToString()
	end, _ := gorest.GetQueryParam(req, "end").ToString()

	log.WithField("trace", "DailySignUp").Infof("Original start and end: %s - %s", start, end)

	period, err := stats.NewPeriod(start, end)
	if err != nil {
		view.Render(w, view.NewBadRequest("Time format must be YYYY-MM-DD"))
		return
	}

	signUps, err := router.model.DailyNewUser(period)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(signUps))
}

// YearlyIncome calculates a year's real income.
//
//	GET /stats/income/year/xxxx
func (router StatsRouter) YearlyIncome(w http.ResponseWriter, req *http.Request) {
	year, err := GetURLParam(req, "year").ToInt()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	y := int(year)
	if y > time.Now().Year() {
		r := view.NewReason()
		r.Field = "year"
		r.Code = "invalid"
		r.SetMessage("Year must be within valid range")
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	fy := subs.NewFiscalYear(y)

	fy, err = router.model.YearlyIncome(fy)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(fy))
}
