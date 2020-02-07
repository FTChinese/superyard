package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/android"
	"gitlab.com/ftchinese/superyard/models/builder"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/apps"
	"net/http"
)

type AndroidRouter struct {
	model apps.AndroidEnv
}

func NewAndroidRouter(db *sqlx.DB) AndroidRouter {
	return AndroidRouter{
		model: apps.AndroidEnv{
			DB: db,
		},
	}
}

// TagExists checks whether a release exists.
func (router AndroidRouter) TagExists(c echo.Context) error {
	versionName := c.Param("versionName")

	ok, err := router.model.Exists(versionName)
	if err != nil {
		return util.NewDBFailure(err)
	}

	if !ok {
		return util.NewNotFound("Version not found")
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateRelease inserts the metadata for a new Android release.
//
// POST /android/releases
//
// Body: {versionName: string, versionCode: int, body: string, apkUrl: string}
func (router AndroidRouter) CreateRelease(c echo.Context) error {
	var r android.Release

	if err := c.Bind(&r); err != nil {
		return util.NewBadRequest(err.Error())
	}

	r.Sanitize()

	if ie := r.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}

	err := router.model.CreateRelease(r)
	if err != nil {
		if util.IsAlreadyExists(err) {
			return util.NewAlreadyExists("versionName")
		}

		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Releases retrieves all releases by sorting version code
// in descending order.
//
// GET /android/releases?page=<number>&per_page=<number>
func (router AndroidRouter) Releases(c echo.Context) error {

	var pagination builder.Pagination
	if err := c.Bind(&pagination); err != nil {
		return util.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	releases, err := router.model.ListReleases(pagination)
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, releases)
}

// SingleReleases retrieves a release by version name
//
// GET /android/releases/{versionName}
func (router AndroidRouter) SingleRelease(c echo.Context) error {
	versionName := c.Param("versionName")

	release, err := router.model.RetrieveRelease(versionName)
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, release)
}

// UpdateRelease updates a single release.
//
// PATCH /android/releases/{versionName}
//
// Body {versionName: string, versionCode: int, body: string, apkUrl: string}
func (router AndroidRouter) UpdateRelease(c echo.Context) error {
	versionName := c.Param("versionName")

	var release android.Release
	if err := c.Bind(&release); err != nil {
		return util.NewBadRequest(err.Error())
	}

	release.Sanitize()

	if ie := release.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}
	release.VersionName = versionName

	if err := router.model.UpdateRelease(release); err != nil {
		if util.IsAlreadyExists(err) {
			return util.NewAlreadyExists("versionCode")
		}

		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteRelease deletes a single release
//
// DELETE /android/releases/:versionName
func (router AndroidRouter) DeleteRelease(c echo.Context) error {
	versionName := c.Param("versionName")

	if err := router.model.DeleteRelease(versionName); err != nil {
		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}
