package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/pkg/android"
	"gitlab.com/ftchinese/superyard/pkg/db"
	"gitlab.com/ftchinese/superyard/repository/apps"
	"net/http"
)

type AndroidRouter struct {
	model    apps.AndroidEnv
	ghClient android.GitHubClient
}

func NewAndroidRouter(db *sqlx.DB) AndroidRouter {
	return AndroidRouter{
		model: apps.AndroidEnv{
			DB: db,
		},
		ghClient: android.MustNewGitHubClient(),
	}
}

// GHLatestRelease get latest release data from GitHub.
// Returns android.Release.
func (router AndroidRouter) GHLatestRelease(c echo.Context) error {
	ghr, respErr := router.ghClient.GetLatestRelease()

	if respErr != nil {
		return respErr
	}

	ghContent, respErr := router.ghClient.GetGradleFile(ghr.TagName)
	if respErr != nil {
		return respErr
	}

	content, err := ghContent.Decode()
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	versionCode, err := android.ParseVersionCode(content)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.JSON(http.StatusOK, ghr.FtcRelease(versionCode))
}

// GHRelease gets a single release from GitHub.
func (router AndroidRouter) GHRelease(c echo.Context) error {
	tag := c.Param("tag")

	ghr, respErr := router.ghClient.GetSingleRelease(tag)

	if respErr != nil {
		return respErr
	}

	ghContent, respErr := router.ghClient.GetGradleFile(ghr.TagName)
	if respErr != nil {
		return respErr
	}

	content, err := ghContent.Decode()
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	versionCode, err := android.ParseVersionCode(content)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.JSON(http.StatusOK, ghr.FtcRelease(versionCode))
}

// TagExists checks whether a release exists in our DB.
// This does not check whether it actually released on GitHub.
func (router AndroidRouter) TagExists(c echo.Context) error {
	versionName := c.Param("versionName")

	ok, err := router.model.Exists(versionName)
	if err != nil {
		return render.NewDBError(err)
	}

	if !ok {
		return render.NewNotFound("")
	}

	return c.NoContent(http.StatusNoContent)
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

	err := router.model.CreateRelease(android.NewRelease(input))
	if err != nil {
		if db.IsAlreadyExists(err) {
			return render.NewAlreadyExists("versionName")
		}

		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Releases retrieves all releases by sorting version code
// in descending order.
//
// GET /android/releases?page=<number>&per_page=<number>
func (router AndroidRouter) Releases(c echo.Context) error {

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

// SingleReleases retrieves a release by version name
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

	err := router.model.UpdateRelease(input)
	if err != nil {
		if db.IsAlreadyExists(err) {
			return render.NewAlreadyExists("versionCode")
		}

		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
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
