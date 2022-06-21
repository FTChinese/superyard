package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AndroidRoutes handlers android app release data.
type AndroidRoutes struct {
	apiClient subsapi.Client
	logger    *zap.Logger
}

// NewAndroidRouter creates a new AndroidRoutes
func NewAndroidRouter(client subsapi.Client, logger *zap.Logger) AndroidRoutes {
	return AndroidRoutes{
		apiClient: client,
		logger:    logger,
	}
}

// ListReleases retrieves all releases by sorting version code
// in descending order.
//
// GET /android/releases?page=<number>&per_page=<number>
func (routes AndroidRoutes) ListReleases(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	claims := getPassportClaims(c)

	rawQuery := c.QueryString()

	resp, err := routes.apiClient.ListAndroidRelease(
		rawQuery,
		claims.Username)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}

// ReleaseOf retrieves a release by version name
//
// GET /android/releases/{versionName}
func (routes AndroidRoutes) ReleaseOf(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	versionName := c.Param("versionName")

	resp, err := routes.apiClient.AndroidReleaseOf(versionName)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}

// CreateRelease inserts the metadata for a new Android release.
//
// POST /android/releases
//
// Body: {versionName: string, versionCode: int, body?: string, apkUrl: string}
func (routes AndroidRoutes) CreateRelease(c echo.Context) error {

	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := routes.
		apiClient.
		CreateAndroidRelease(
			c.Request().Body,
			claims.Username)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}

// UpdateRelease updates a single release.
//
// PATCH /android/releases/{versionName}
//
// Body {body: string, apkUrl: string}
func (routes AndroidRoutes) UpdateRelease(c echo.Context) error {
	versionName := c.Param("versionName")

	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := routes.
		apiClient.
		UpdateAndroidRelease(
			versionName,
			c.Request().Body,
			claims.Username)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}

// DeleteRelease deletes a single release
//
// DELETE /android/releases/:versionName
func (routes AndroidRoutes) DeleteRelease(c echo.Context) error {
	versionName := c.Param("versionName")

	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := routes.
		apiClient.
		DeleteAndroidRelease(
			versionName,
			claims.Username)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}
