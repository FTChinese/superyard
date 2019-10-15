package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/oauth"
	"gitlab.com/ftchinese/backyard-api/repository/registry"
	"gitlab.com/ftchinese/backyard-api/repository/staff"
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
func (router ApiRouter) CreateApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var app oauth.App
	if err := gorest.ParseJSON(req.Body, &app); err != nil {
		_ = view.Render(w, view.NewBadRequest(""))
		return
	}

	app.Sanitize()

	logger.WithField("trace", "CreateApp").Infof("%+v", app)

	if r := app.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	if err := app.GenCredentials(); err != nil {
		_ = view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	app.OwnedBy = userName

	err := router.model.CreateApp(app)

	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Code = view.CodeAlreadyExists
			reason.Field = "slug"
			_ = view.Render(w, view.NewUnprocessable(reason))
			return
		}

		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// ListApps loads all app with pagination support
//
//	GET /next/apps?page=<number>&per_page=<number>
func (router ApiRouter) ListApps(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	apps, err := router.model.ListApps(pagination)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(apps))
}

// RetrieveApp retrieves an app by its slug name.
func (router ApiRouter) LoadApp(w http.ResponseWriter, req *http.Request) {
	clientID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	app, err := router.model.RetrieveApp(clientID)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(app))
}

// UpdateApp updates an app's data.
//
//	PATCH /api/apps/{id}
//
// Input {name: string, slug: string, repoUrl?: string, description?: string, homeUrl?: string, callbackUrl?: string, ownedBy: string}
func (router ApiRouter) UpdateApp(w http.ResponseWriter, req *http.Request) {

	clientID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var app oauth.App
	if err := gorest.ParseJSON(req.Body, &app); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	app.Sanitize()
	if r := app.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	app.ClientID = clientID

	if err := router.model.UpdateApp(app); err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Code = view.CodeAlreadyExists
			reason.Field = "slug"
			_ = view.Render(w, view.NewUnprocessable(reason))
			return
		}
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// RemoveApp flags an app as inactive.
// This also removes all access tokens owned by this app.
//
//	DELETE /api/apps/{id}
// Input: {ownedBy: string}
func (router ApiRouter) RemoveApp(w http.ResponseWriter, req *http.Request) {

	clientID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var by oauth.AppRemover
	if err := gorest.ParseJSON(req.Body, &by); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	by.ClientID = clientID

	err = router.model.RemoveApp(by)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// ListKeys shows all access tokens owned by an app or by a human.
func (router ApiRouter) ListKeys(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var by oauth.KeySelector
	if err := decoder.Decode(&by, req.Form); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p := gorest.GetPagination(req)

	tokens, err := router.model.ListKeys(by, p)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(tokens))
}

// NewToken creates an access token for a person or for an app.
// Input: {description: string, createdBy: string, clientId?: string}
func (router ApiRouter) CreateKey(w http.ResponseWriter, req *http.Request) {

	var input oauth.InputKey
	if err := gorest.ParseJSON(req.Body, &input); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	acc, err := oauth.NewAccess(input)
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	_, err = router.model.CreateToken(acc)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// DeletePersonalKeys deactivates all keys owned by a user.
// Input: { createdBy: string}
func (router ApiRouter) DeletePersonalKeys(w http.ResponseWriter, req *http.Request) {
	var by oauth.KeyRemover
	if err := gorest.ParseJSON(req.Body, &by); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if by.CreatedBy == "" {
		r := view.NewReason()
		r.Field = "createdBy"
		r.Code = view.CodeMissingField
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	if err := router.model.DeleteKeys(by.CreatedBy); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// RemoveKey deactivate an access token created by a user.
// The token could be owned by either an app or a human being.
func (router ApiRouter) RemoveKey(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToInt()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var by oauth.KeyRemover
	if err := gorest.ParseJSON(req.Body, &by); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	by.ID = id

	if err := router.model.RemoveKey(by); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}
