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
		view.Render(w, view.NewBadRequest(""))
		return
	}

	app.Sanitize()

	logger.WithField("trace", "CreateApp").Infof("%+v", app)

	if r := app.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	err := app.GenCredentials()
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	app.OwnedBy = userName

	err = router.model.CreateApp(app)

	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Code = view.CodeAlreadyExists
			reason.Field = "slug"
			view.Render(w, view.NewUnprocessable(reason))
			return
		}

		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// ListApps loads all app with pagination support
//
//	GET /next/apps?page=<number>&per_page=<number>
func (router ApiRouter) ListApps(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	apps, err := router.model.ListApps(pagination)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(apps))
}

// RetrieveApp retrieves an app by its slug name.
//
// Get /next/apps/{name}
func (router ApiRouter) LoadApp(w http.ResponseWriter, req *http.Request) {
	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	app, err := router.model.RetrieveApp(slugName)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(app))
}

// UpdateApp updates an app's data.
//
//	PATCH /next/apps/{name}
//
// Input {name: string, slug: string, repoUrl: string, description: string, homeUrl: string}
func (router ApiRouter) UpdateApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	//slugName, err := GetURLParam(req, "name").ToString()
	//if err != nil {
	//	view.Render(w, view.NewBadRequest(err.Error()))
	//	return
	//}

	var app oauth.App
	if err := gorest.ParseJSON(req.Body, &app); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	app.Sanitize()
	if r := app.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	app.OwnedBy = userName

	err := router.model.UpdateApp(app)
	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Code = view.CodeAlreadyExists
			reason.Field = "slug"
			view.Render(w, view.NewUnprocessable(reason))
			return
		}
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// RemoveApp flags an app as inactive.
// This also removes all access tokens owned by this app.
//
//	DELETE /next/apps/{name}
func (router ApiRouter) RemoveApp(w http.ResponseWriter, req *http.Request) {

	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	clientID, err := router.model.SearchApp(slugName)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	err = router.model.RemoveApp(clientID)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// NewToken creates an access token for a person or for an app.
//
//	POST /next/apps/{name}/tokens
//
// Input: {description: string}
func (router ApiRouter) CreateKey(w http.ResponseWriter, req *http.Request) {

	acc, err := oauth.NewAccess()
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}
	if err := gorest.ParseJSON(req.Body, &acc); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	_, err = router.model.CreateToken(acc)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

func (router ApiRouter) ListKeys(w http.ResponseWriter, req *http.Request) {

}

func (router ApiRouter) RemoveKey(w http.ResponseWriter, req *http.Request) {

}
