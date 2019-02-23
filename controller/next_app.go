package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/model"
	"gitlab.com/ftchinese/backyard-api/oauth"
	"net/http"
)

type NextAPIRouter struct {
	model model.OAuthEnv
	staff model.StaffEnv
}

// NewNextAPIRouter creates a new instance of FTCAPIRouter.
func NewNextAPIRouter(db *sql.DB) NextAPIRouter {
	return NextAPIRouter{
		model: model.OAuthEnv{DB: db},
		staff: model.StaffEnv{DB: db},
	}
}

// CreateApp creates an new app which needs to access next-api.
//
//	POST /next/apps
//
// Input {name: string, slug: string, repoUrl: string, description: string, homeUrl: string}
func (router NextAPIRouter) CreateApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var app oauth.App
	if err := gorest.ParseJSON(req.Body, &app); err != nil {
		view.Render(w, view.NewBadRequest(""))
		return
	}

	app.Sanitize()

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

	err = router.model.SaveApp(app)

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
func (router NextAPIRouter) ListApps(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := GetPagination(req)

	apps, err := router.model.ListApps(pagination)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(apps))
}

// LoadApp retrieves an app by its slug name.
//
// Get /next/apps/{name}
func (router NextAPIRouter) LoadApp(w http.ResponseWriter, req *http.Request) {
	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	app, err := router.model.LoadApp(slugName)
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
func (router NextAPIRouter) UpdateApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

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

	err = router.model.UpdateApp(slugName, app)
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
func (router NextAPIRouter) RemoveApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	clientID, err := router.model.FindClientID(slugName)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	err = router.model.RemoveApp(clientID, userName)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// TransferApp changes ownership of an app
//
//	POST /next/apps/{name}/transfer
//
// Input {newOwner: string}
func (router NextAPIRouter) TransferApp(w http.ResponseWriter, req *http.Request) {
	currentUser := req.Header.Get(userNameKey)

	slugName, err := GetURLParam(req, "name").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var o oauth.Ownership
	if err := gorest.ParseJSON(req.Body, &o); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}
	o.SlugName = slugName
	o.OldOwner = currentUser

	// TODO: validate and sanitize.

	exists, err := router.staff.NameExists(o.NewOwner)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}
	if !exists {
		view.Render(w, view.NewNotFound())
		return
	}

	err = router.model.TransferApp(o)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}
