package controller

import (
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/ftchinese/backyard-api/ftcapi"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// FTCController creates routers to manipulate ftc apps and api keys
// All routers requires `X-User-Name` header
type FTCController struct {
	ftcModel   ftcapi.Env
	staffModel staff.Env
}

// NewApp creates an new app built on ftc api
// Input:
// {
//	name: string,
//	slug: string,
//	repoUrl: string,
//	description: string,
//	homeUrl: string
// }
func (c FTCController) NewApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var app ftcapi.App
	// 400 Bad Request
	// { message: "Problems parsing JSON" }
	if err := util.Parse(req.Body, &app); err != nil {
		view.Render(w, util.NewBadRequest(""))
		return
	}

	app.Sanitize()

	// TODO: validation

	app.OwnedBy = userName

	err := c.ftcModel.NewApp(app)

	// { message: "Validation failed",
	// 	field: "slug",
	//	code: "already_exists"
	// }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "slug"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// ListApps loads all app with pagination support
// TODO: add a middleware to parse form.
func (c FTCController) ListApps(w http.ResponseWriter, req *http.Request) {
	page, err := getQueryParam(req, "page").toInt()

	if err != nil {
		page = 1
	}

	apps, err := c.ftcModel.AppRoster(page, 20)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(apps))
}

// GetApp loads an app of the specified slug name
func (c FTCController) GetApp(w http.ResponseWriter, req *http.Request) {
	slugName := chi.URLParam(req, "name")

	// 400 Bad Request
	// { message: "Invalid request URI" }
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	app, err := c.ftcModel.RetrieveApp(slugName)

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(app))
}

// UpdateApp updates an app's data
// Input:
// {
//	name: string,
//	slug: string,
//	repoUrl: string,
//	description: string,
//	homeUrl: string
// }
func (c FTCController) UpdateApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	slugName := chi.URLParam(req, "name")

	// 400 Bad Request
	// { message: "Invalid request URI" }
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	var app ftcapi.App
	// 400 Bad Request
	// { message: "Problems parsing JSON" }
	if err := util.Parse(req.Body, &app); err != nil {
		view.Render(w, util.NewBadRequest(""))
		return
	}

	app.Sanitize()

	// TODO: validation

	// OwnedBy is used to make sure the update operaton is performed by the owner
	app.OwnedBy = userName

	err := c.ftcModel.UpdateApp(slugName, app)

	// { message: "Validation failed",
	// 	field: "slug",
	//	code: "already_exists"
	// }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "slug"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// RemoveApp flags an app as inactive
// This also removes all access tokens owned by this app
func (c FTCController) RemoveApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	slugName := chi.URLParam(req, "name")

	// 400 Bad Request
	// { message: "Invalid request URI" }
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := c.ftcModel.RemoveApp(slugName, userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, "slug"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// TransferApp changes ownership of an app
// Input {newOwner: string}
func (c FTCController) TransferApp(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	slugName := chi.URLParam(req, "name")

	// 400 Bad Request
	// { message: "Invalid request URI" }
	if slugName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	newOwner, err := util.GetJSONString(req.Body, "newOwner")

	// 400 Bad Request
	// { message: "Problems parsing JSON" }
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// First check if target owner exists
	exists, err := c.staffModel.StaffNameExists(newOwner)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
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
	err = c.ftcModel.TransferApp(o)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}
