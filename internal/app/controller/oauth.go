package controller

import (
	"net/http"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/internal/app/repository/registry"
	"github.com/FTChinese/superyard/internal/pkg/oauth"
	"github.com/FTChinese/superyard/pkg/conv"
	"github.com/FTChinese/superyard/pkg/db"

	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
)

type OAuthRouter struct {
	regRepo registry.Env
}

// NewOAuthRouter creates a new instance of FTCAPIRouter.
func NewOAuthRouter(myDBs db.ReadWriteMyDBs) OAuthRouter {
	return OAuthRouter{
		regRepo: registry.NewEnv(myDBs),
	}
}

// CreateApp creates an new app which needs to access next-api.
//
// Input {name: string, slug: string, repoUrl: string, description: string, homeUrl: string}
func (router OAuthRouter) CreateApp(c echo.Context) error {
	claims := getPassportClaims(c)

	var input oauth.BaseApp
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	input.Sanitize()

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	app, err := oauth.NewApp(input, claims.Username)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}
	app.OwnedBy = claims.Username

	err = router.regRepo.CreateApp(app)
	if err != nil {
		if db.IsAlreadyExists(err) {
			return render.NewAlreadyExists("slug")
		}
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, app)
}

// ListApps loads all app with pagination support
//
//	GET /next/apps?page=<number>&per_page=<number>
func (router OAuthRouter) ListApps(c echo.Context) error {

	var pagination gorest.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	apps, err := router.regRepo.ListApps(pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, apps)
}

// RetrieveApp retrieves an app by its client id.
func (router OAuthRouter) LoadApp(c echo.Context) error {
	clientID := c.Param("id")

	app, err := router.regRepo.RetrieveApp(clientID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, app)
}

// UpdateApp updates an app's data.
// Any logged in user could update any app, regardless who owned it.
//
//	PATCH /api/apps/:id
//
// Input {name: string, slug: string, repoUrl?: string, description?: string, homeUrl?: string, callbackUrl?: string}
func (router OAuthRouter) UpdateApp(c echo.Context) error {

	clientID := c.Param("id")

	var input oauth.BaseApp
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	input.Sanitize()
	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	app := oauth.App{
		BaseApp:  input,
		ClientID: clientID,
	}

	if err := router.regRepo.UpdateApp(app); err != nil {
		if db.IsAlreadyExists(err) {
			return render.NewAlreadyExists("slug")
		}
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveApp flags an app as inactive.
// This also removes all access tokens owned by this app.
// We does not impose access control here.
// Anyone can delete any app created by others.
//
//	DELETE /api/apps/:id
func (router OAuthRouter) RemoveApp(c echo.Context) error {
	clientID := c.Param("id")

	if err := router.regRepo.RemoveApp(clientID); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListKeys shows all access tokens owned by an app or by a human.
// Query params:
// ?client_id=<string>&page=<number>&per_page=<number>
// All are optional.
// Is client_id is present, it indicates the client
// is requesting tokens belong to an app; otherwise it
// indicates personal keys.
func (router OAuthRouter) ListKeys(c echo.Context) error {

	clientID := c.QueryParam("client_id")
	claims := getPassportClaims(c)

	var tokens []oauth.Access
	var err error
	if clientID != "" {
		tokens, err = router.regRepo.ListAppTokens(clientID)
	} else {
		tokens, err = router.regRepo.ListPersonalKeys(claims.Username)
	}

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, tokens)
}

// CreateKey creates an access token for a person or for an app.
// Input: {description?: string, clientId?: string}
func (router OAuthRouter) CreateKey(c echo.Context) error {
	claims := getPassportClaims(c)

	var input oauth.BaseAccess
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	acc, err := oauth.NewAccess(input, claims.Username)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	_, err = router.regRepo.CreateToken(acc)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, acc)
}

// RemoveKey deactivate an access token created by a user.
// The token could be owned by either an app or a human being.
func (router OAuthRouter) RemoveKey(c echo.Context) error {
	claims := getPassportClaims(c)

	id, err := conv.ParseInt64(c.Param("id"))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	key := oauth.Access{
		ID:        id,
		CreatedBy: claims.Username,
	}

	if err := router.regRepo.RemoveKey(key); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
