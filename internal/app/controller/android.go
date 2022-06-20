package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/apps"
	"github.com/FTChinese/superyard/internal/pkg/android"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AndroidRouter struct {
	model apps.Env
}

func NewAndroidRouter(myDBs db.ReadWriteMyDBs) AndroidRouter {
	return AndroidRouter{
		model: apps.NewEnv(myDBs),
	}
}

// CreateRelease inserts the metadata for a new Android release.
//
// POST /android/releases
//
// Body: {versionName: string, versionCode: int, body?: string, apkUrl: string}
func (router AndroidRouter) CreateRelease(c echo.Context) error {
	var input android.ReleaseInput

	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateCreation(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	release, err := router.model.CreateRelease(android.NewRelease(input))
	if err != nil {
		if db.IsAlreadyExists(err) {
			return render.NewAlreadyExists("versionName")
		}

		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, release)
}

// ListReleases retrieves all releases by sorting version code
// in descending order.
//
// GET /android/releases?page=<number>&per_page=<number>
func (router AndroidRouter) ListReleases(c echo.Context) error {

	var pagination gorest.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	releases, err := router.model.ListReleases(pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, releases)
}

// SingleRelease retrieves a release by version name
//
// GET /android/releases/{versionName}
func (router AndroidRouter) SingleRelease(c echo.Context) error {
	versionName := c.Param("versionName")

	release, err := router.model.RetrieveRelease(versionName)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, release)
}

// UpdateRelease updates a single release.
//
// PATCH /android/releases/{versionName}
//
// Body {body: string, apkUrl: string}
func (router AndroidRouter) UpdateRelease(c echo.Context) error {
	versionName := c.Param("versionName")

	var input android.ReleaseInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	input.VersionName = versionName

	if ve := input.ValidateUpdate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	current, err := router.model.RetrieveRelease(versionName)
	if err != nil {
		return render.NewDBError(err)
	}

	updated := current.Update(input)

	err = router.model.UpdateRelease(updated)
	if err != nil {
		if db.IsAlreadyExists(err) {
			return render.NewAlreadyExists("versionCode")
		}

		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, updated)
}

// DeleteRelease deletes a single release
//
// DELETE /android/releases/:versionName
func (router AndroidRouter) DeleteRelease(c echo.Context) error {
	versionName := c.Param("versionName")

	if err := router.model.DeleteRelease(versionName); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
