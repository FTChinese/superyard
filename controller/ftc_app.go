package controller

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/ftcapi"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// FTCAPIRouter controls access to next-api by human and applications.
// All routers requires `X-User-Name` header.
type FTCAPIRouter struct {
	apiModel   ftcapi.Env
	staffModel staff.Env // used to check if a staff exists
}

// NewFTCAPIRouter creates a new instance of FTCAPIRouter.
func NewFTCAPIRouter(db *sql.DB) FTCAPIRouter {
	api := ftcapi.Env{DB: db}
	staff := staff.Env{DB: db}

	return FTCAPIRouter{
		apiModel:   api,
		staffModel: staff,
	}
}

// NewApp creates an new app which needs to access next-api.
//
//	POST /ftc-api/apps
func (c FTCAPIRouter) NewApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var app ftcapi.App

	// 400 Bad Request
	if err := util.Parse(req.Body, &app); err != nil {
		view.Render(w, util.NewBadRequest(""))
		return
	}

	app.Sanitize()

	// 422 Unprocessable Entity
	if r := app.Validate(); r != nil {
		view.Render(w, util.NewUnprocessable(r))
		return
	}

	app.OwnedBy = userName

	err := c.apiModel.NewApp(app)

	// Duplicate error
	if err != nil {
		if util.IsAlreadyExists(err) {
			reason := util.NewReasonAlreadyExists("slug")

			view.Render(w, util.NewUnprocessable(reason))

			return
		}

		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewNoContent())
}

// ListApps loads all app with pagination support
//
//	GET /ftc-api/apps?page=<number>
func (c FTCAPIRouter) ListApps(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request if query string cannot be parsed.
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))
		return
	}

	page, err := getQueryParam(req, "page").toInt()

	if err != nil {
		page = 1
	}

	apps, err := c.apiModel.AppRoster(page, 20)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, util.NewResponse().NoCache().SetBody(apps))
}

// GetApp loads an app.
//
//	GET /ftc-api/apps/{name}
func (c FTCAPIRouter) GetApp(w http.ResponseWriter, req *http.Request) {
	slugName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	app, err := c.apiModel.RetrieveApp(slugName)

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, util.NewResponse().NoCache().SetBody(app))
}

// UpdateApp updates an app's data.
//
//	PATCH /ftc-api/apps/{name}
func (c FTCAPIRouter) UpdateApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	slugName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	var app ftcapi.App
	// 400 Bad Request
	if err := util.Parse(req.Body, &app); err != nil {
		view.Render(w, util.NewBadRequest(""))
		return
	}

	app.Sanitize()

	// 422 Unprocessable Entity
	if r := app.Validate(); r != nil {
		view.Render(w, util.NewUnprocessable(r))
		return
	}

	// OwnedBy is used to make sure the update operaton is performed by the owner
	app.OwnedBy = userName

	err := c.apiModel.UpdateApp(slugName, app)

	// 422 Unprocessable Entity
	if err != nil {
		if util.IsAlreadyExists(err) {
			reason := util.NewReasonAlreadyExists("slug")
			view.Render(w, util.NewUnprocessable(reason))
			return
		}
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewNoContent())
}

// DeleteApp flags an app as inactive.
// This also removes all access tokens owned by this app.
//
//	DELETE /ftc-api/apps/{name}
func (c FTCAPIRouter) DeleteApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	slugName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := c.apiModel.RemoveApp(slugName, userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}

// TransferApp changes ownership of an app
//
//	POST /ftc-api/apps/{name}/transfer
func (c FTCAPIRouter) TransferApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	slugName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	newOwner, err := util.GetJSONString(req.Body, "newOwner")

	// 400 Bad Request
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// First check if target owner exists
	exists, err := c.staffModel.StaffNameExists(newOwner)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	// 404 Not Found
	if !exists {
		view.Render(w, util.NewNotFound())

		return
	}

	o := ftcapi.Ownership{
		SlugName: slugName,
		NewOwner: newOwner,
		OldOwner: userName,
	}
	err = c.apiModel.TransferApp(o)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}
