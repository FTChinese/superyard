package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/stst"
	"github.com/FTChinese/superyard/internal/pkg/stats"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// StatsRouter responds to requests for statistic data.
type StatsRouter struct {
	repo stst.Env
}

// NewStatsRouter creates a new instance of StatsRouter
func NewStatsRouter(myDBs db.ReadWriteMyDBs) StatsRouter {

	return StatsRouter{
		repo: stst.NewEnv(myDBs),
	}
}

func (router StatsRouter) AliUnconfirmed(c echo.Context) error {
	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}
	page.Normalize()

	unconfirmed, err := router.repo.AliUnconfirmed(page)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, unconfirmed)
}

func (router StatsRouter) WxUnconfirmed(c echo.Context) error {
	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}
	page.Normalize()

	unconfirmed, err := router.repo.WxUnconfirmed(page)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, unconfirmed)
}

// DailySignUp show how many new users signed up at ftchinese.com everyday.
//
//	GET /stats/signup/daily?start=YYYY-MM-DD&end=YYYY-MM-DD
func (router StatsRouter) DailySignUp(c echo.Context) error {

	start := c.QueryParam("start")
	end := c.QueryParam("end")

	period, err := stats.NewPeriod(start, end)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	signUps, err := router.repo.DailyNewUser(period)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, signUps)
}

// YearlyIncome calculates a year's real income.
//
//	GET /stats/income/year/xxxx
func (router StatsRouter) YearlyIncome(c echo.Context) error {
	year, err := conv.ParseInt64(c.Param("year"))
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

	fy := stats.NewFiscalYear(y)

	fy, err = router.repo.YearlyIncome(fy)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, fy)
}
