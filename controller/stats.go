package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/promo"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/models/validator"
	"gitlab.com/ftchinese/superyard/repository/aggregate"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/superyard/models/stats"
)

// StatsRouter responds to requests for statistic data.
type StatsRouter struct {
	model aggregate.StatsEnv
}

// NewStatsRouter creates a new instance of StatsRouter
func NewStatsRouter(db *sqlx.DB) StatsRouter {

	return StatsRouter{
		model: aggregate.StatsEnv{DB: db},
	}
}

// DailySignUp show how many new users signed up at ftchinese.com everyday.
//
//	GET /stats/signup/daily?start=YYYY-MM-DD&end=YYYY-MM-DD
func (router StatsRouter) DailySignUp(c echo.Context) error {

	start := c.QueryParam("start")
	end := c.QueryParam("end")

	log.WithField("trace", "DailySignUp").Infof("Original start and end: %s - %s", start, end)

	period, err := stats.NewPeriod(start, end)
	if err != nil {
		return util.NewBadRequest(err.Error())
	}

	signUps, err := router.model.DailyNewUser(period)

	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, signUps)
}

// YearlyIncome calculates a year's real income.
//
//	GET /stats/income/year/xxxx
func (router StatsRouter) YearlyIncome(c echo.Context) error {
	year, err := ParseInt(c.Param("year"))
	if err != nil {
		return util.NewBadRequest(err.Error())
	}

	y := int(year)
	if y > time.Now().Year() {
		return util.NewUnprocessable(&validator.InputError{
			Message: "Year must be within valid range",
			Field:   "year",
			Code:    validator.CodeInvalid,
		})
	}

	fy := promo.NewFiscalYear(y)

	fy, err = router.model.YearlyIncome(fy)

	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, fy)
}
