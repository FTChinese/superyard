package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/oauth"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/models/validator"
	"gitlab.com/ftchinese/superyard/repository/registry"
	"gitlab.com/ftchinese/superyard/repository/staff"
	"net/http"
)

type ApiRouter struct {
	model registry.Env
	staff staff.Env
}

// APIRouter creates a new instance of FTCAPIRouter.
func APIRouter(db *sqlx.DB) ApiRouter {
	return ApiRouter{
		model: registry.Env{DB: db},
		staff: staff.Env{DB: db},
	}
}

// CreateApp creates an new app which needs to access next-api.
//
//	POST /next/apps
//
// Input {name: string, slug: string, repoUrl: string, description: string, homeUrl: string}
func (router ApiRouter) CreateApp(c echo.Context) error {
	userName := c.Request().Header.Get(userNameKey)

	var app oauth.App
	if err := c.Bind(&app); err != nil {
		return render.NewBadRequest(err.Error())
	}

	app.Sanitize()

	logger.WithField("trace", "CreateApp").Infof("%+v", app)

	if ve := app.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	if err := app.GenCredentials(); err != nil {
		return render.NewBadRequest(err.Error())
	}

	app.OwnedBy = userName

	err := router.model.CreateApp(app)

	if err != nil {
		if util.IsAlreadyExists(err) {
			return render.NewAlreadyExists("slug")
		}

		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListApps loads all app with pagination support
//
//	GET /next/apps?page=<number>&per_page=<number>
func (router ApiRouter) ListApps(c echo.Context) error {

	var pagination util.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	apps, err := router.model.ListApps(pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, apps)
}

// RetrieveApp retrieves an app by its slug name.
func (router ApiRouter) LoadApp(c echo.Context) error {
	clientID := c.Param("id")

	app, err := router.model.RetrieveApp(clientID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, app)
}

// UpdateApp updates an app's data.
//
//	PATCH /api/apps/{id}
//
// Input {name: string, slug: string, repoUrl?: string, description?: string, homeUrl?: string, callbackUrl?: string, ownedBy: string}
func (router ApiRouter) UpdateApp(c echo.Context) error {

	clientID := c.Param("id")

	var app oauth.App
	if err := c.Bind(&app); err != nil {
		return render.NewBadRequest(err.Error())
	}

	app.Sanitize()
	if ve := app.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	app.ClientID = clientID

	if err := router.model.UpdateApp(app); err != nil {
		if util.IsAlreadyExists(err) {
			return render.NewAlreadyExists("slug")
		}
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveApp flags an app as inactive.
// This also removes all access tokens owned by this app.
//
//	DELETE /api/apps/{id}
// Input: {ownedBy: string}
func (router ApiRouter) RemoveApp(c echo.Context) error {

	clientID := c.Param("id")

	var by oauth.AppRemover
	if err := c.Bind(&by); err != nil {
		return render.NewBadRequest(err.Error())
	}

	by.ClientID = clientID

	if err := router.model.RemoveApp(by); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListKeys shows all access tokens owned by an app or by a human.
func (router ApiRouter) ListKeys(c echo.Context) error {

	var selector oauth.KeySelector
	if err := c.Bind(&selector); err != nil {
		return render.NewBadRequest(err.Error())
	}

	var p util.Pagination
	if err := c.Bind(&p); err != nil {
		return render.NewBadRequest(err.Error())
	}
	p.Normalize()

	tokens, err := router.model.ListKeys(selector, p)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, tokens)
}

// NewToken creates an access token for a person or for an app.
// Input: {description?: string, createdBy: string, clientId?: string}
func (router ApiRouter) CreateKey(c echo.Context) error {

	var input oauth.InputKey
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	acc, err := oauth.NewAccess(input)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	_, err = router.model.CreateToken(acc)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeletePersonalKeys deactivates all keys owned by a user.
// Input: { createdBy: string}
func (router ApiRouter) DeletePersonalKeys(c echo.Context) error {
	var by oauth.KeyRemover
	if err := c.Bind(&by); err != nil {
		return render.NewBadRequest(err.Error())
	}

	ve := validator.New("createdBy").Required().Validate(by.CreatedBy)
	if ve != nil {
		return render.NewUnprocessable(ve)
	}

	if err := router.model.DeleteKeys(by.CreatedBy); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveKey deactivate an access token created by a user.
// The token could be owned by either an app or a human being.
// Input { createdBy: string}
// The input restricts that user could only delete keys created by itself.
func (router ApiRouter) RemoveKey(c echo.Context) error {
	id, err := ParseInt(c.Param("id"))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	var by oauth.KeyRemover
	if err := c.Bind(&by); err != nil {
		return render.NewBadRequest(err.Error())
	}

	by.ID = id

	if err := router.model.RemoveKey(by); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
