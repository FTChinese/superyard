package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/push"
	apn2 "github.com/FTChinese/superyard/repository/apn"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

type APNRouter struct {
	model apn2.APNEnv
}

func NewAPNRouter(db *sqlx.DB) APNRouter {
	return APNRouter{
		model: apn2.APNEnv{DB: db},
	}
}

func (router APNRouter) ListMessages(c echo.Context) error {
	var pagination gorest.Pagination
	// 400 Bad Request if query string cannot be parsed.
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	msgs, err := router.model.ListMessage(pagination)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, msgs)
}

func (router APNRouter) LoadTimezones(c echo.Context) error {
	tz, err := router.model.TimeZoneDist()

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, tz)
}

func (router APNRouter) LoadDeviceDist(c echo.Context) error {
	d, err := router.model.DeviceDist()

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, d)
}

func (router APNRouter) LoadInvalidDist(c echo.Context) error {
	d, err := router.model.InvalidDist()

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, d)
}

func (router APNRouter) CreateTestDevice(c echo.Context) error {
	var d push.TestDevice

	if err := c.Bind(&d); err != nil {
		return render.NewBadRequest(err.Error())
	}

	err := router.model.CreateTestDevice(d)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router APNRouter) ListTestDevice(c echo.Context) error {
	d, err := router.model.ListTestDevice()

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, d)
}

func (router APNRouter) RemoveTestDevice(c echo.Context) error {
	id, err := ParseInt(c.Param("id"))

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	err = router.model.RemoveTestDevice(id)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
