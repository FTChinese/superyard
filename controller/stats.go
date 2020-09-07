package controller

import (
	"github.com/FTChinese/go-rest/render"
	stats2 "github.com/FTChinese/superyard/pkg/stats"
	"github.com/FTChinese/superyard/repository/aggregate"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
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

	period, err := stats2.NewPeriod(start, end)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	signUps, err := router.model.DailyNewUser(period)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, signUps)
}

// YearlyIncome calculates a year's real income.
//
//	GET /stats/income/year/xxxx
func (router StatsRouter) YearlyIncome(c echo.Context) error {
	year, err := ParseInt(c.Param("year"))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	y := int(year)
	if y > time.Now().Year() {
		return render.NewUnprocessable(&render.ValidationError{
			Message: "Year must be within valid range",
			Field:   "year",
			Code:    render.CodeInvalid,
		})
	}

	fy := stats2.NewFiscalYear(y)

	fy, err = router.model.YearlyIncome(fy)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, fy)
}
