package controller

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/ftcapi"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// FTCAPIRouter controls access to next-api.
// All routers requires `X-User-Name` header.
//
// * POST `/ftc-api/apps` Create a new app that needs to access next-api.
//
// * GET `/ftc-api/apps?page=<number>` Show all ftc apps. Anyone can see details of an app created by any others.
//
// * GET `/ftc-api/apps/{name}` Show the detial of a ftc app
//
// * PATCH `/ftc-api/apps/{name}` Allow owner of an app to edit it.
//
// * DELETE `/ftc-api/apps/{name}` Delete an app.
//
// * POST `/ftc-api/apps/{name}/transfer` Transfer ownership of an app to others.
//
// * POST `/ftc-api/tokens` Create an access token. It could belong to a person or an app, depending on the data passed in.
//
// * GET `/ftc-api/tokens/personal` Show all access tokens granted to a user.
//
// * DELETE `/ftc-api/token/personal/{tokenId}` Revoke an access token owned by a user.
//
// * GET `/ftc-api/tokens/app/{name}` Show all access tokens owned by an app.
//
// * DELETE `/ftc-api/tokens/app/{name}/{tokenId}` Revoke an access token owned by an app.
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
//
// Input:
// 	{
//		"name": "User Login", // required, max 255 chars
//		"slug": "user-login", // required, max 255 chars
//		"repoUrl": "https://github.com/user-login", // required, 120 chars
//		"description": "UI for user login", // optional, 511 chars
//		"homeUrl": "https://www.ftchinese.com/user" // optional, 255 chars
// 	}
//
// - `400 Bad Request` if request body cannot be parsed as JSON.
//	{
// 		"message": "Problems parsing JSON"
//	}
//
// - 422 Unprocessable Entity if required fields are missing,
// 	{
// 		"message": "Validation failed",
// 		"field": "name | slug | repoUrl",
// 		"code": "missing"
// 	}
// or the length of  any of the fields exceeds max chars,
// 	{
// 		"message": "The length of xxx should not exceed 255 chars",
// 		"field": "email | slug | repoUrl | description | homeUrl",
// 		"code": "invalid"
// 	}
// or the slugified name of the app is taken
//	{
//		"message": "Validation failed",
// 		"field": "slug",
//		"code": "already_exists"
// 	}
//
// - `204 No Content` for success.
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
	if r := app.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))
		return
	}

	app.OwnedBy = userName

	err := c.apiModel.NewApp(app)

	// Duplicate error
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "slug"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// ListApps loads all app with pagination support
//
//	GET /ftc-api/apps?page=<number>
//
// `page` defaults to 1 if it is missing, or is not a number.
//
// - 400 Bad Request if query string cannot be parsed.
//
// - 200 OK with body:
// 	[{
//		"name": "User Login",
//		"slug": "user-login",
//		"clientId": "20 hexdecimal numbers"
// 		"clientSecret": "64 hexdecimal numbers"
//		"repoUrl": "https://github.com/user-login",
//		"description": "UI for user login",
//		"homeUrl": "https://www.ftchinese.com/user",
// 		"isActive": true,
// 		"createdAt": "",
// 		"updatedAt": "",
// 		"ownedBy": "foo.bar"
// }]
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
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// 200 OK
	view.Render(w, util.NewResponse().NoCache().SetBody(apps))
}

// GetApp loads an app.
//
//	GET /ftc-api/apps/{name}
//
// - `400 Bad Request` if request URL does not contain `name` part
//	{
//		"message": "Invalid request URI"
//	}
//
// - `404 Not Found` if the app does not exist
//
// - 200 OK. See response for ListApps.
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
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// 200 OK
	view.Render(w, util.NewResponse().NoCache().SetBody(app))
}

// UpdateApp updates an app's data.
//
//	PATCH /ftc-api/apps/{name}
//
// Input:
// 	{
//		"name": "User Login", // max 60 chars, required
//		"slug": "user-login", // max 60 chars, required
//		"repoUrl": "https://github.com/user-login", // 120 chars, required
//		"description": "UI for user login", // 500 chars, optional
//		"homeUrl": "https://www.ftchinese.com/user" // 120 chars, optional
// }
//
// - `400 Bad Request` if request URL does not contain `name` part
//	{
//		"message": "Invalid request URI"
//	}
// or if request body cannot be parsed as JSON.
//	{
// 		"message": "Problems parsing JSON"
//	}
//
// - 422 Unprocessable Entity is the same as `POST /ftc-api/apps` used by NewApp()
//
// - `204 No Content` for success.
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
	if r := app.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))
		return
	}

	// OwnedBy is used to make sure the update operaton is performed by the owner
	app.OwnedBy = userName

	err := c.apiModel.UpdateApp(slugName, app)

	// 422 Unprocessable Entity
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "slug"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// DeleteApp flags an app as inactive.
// This also removes all access tokens owned by this app.
//
//	DELETE /ftc-api/apps/{name}
//
// - `400 Bad Request` if request URL does not contain `name` part
//	{
//		"message": "Invalid request URI"
//	}
//
// - `204 No Content` for success.
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
		view.Render(w, util.NewDBFailure(err, "slug"))

		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}

// TransferApp changes ownership of an app
//
//	POST /ftc-api/apps/{name}/transfer
//
// Input
// 	{
// 		"newOwner": "foo.baz"
// 	}
//
// - `400 Bad Request` if request URL does not contain `name` part
//	{
//		"message": "Invalid request URI"
//	}
// or if request body cannot be parsed as JSON.
//	{
// 		"message": "Problems parsing JSON"
//	}
//
// - 404 Not Found if the new owner is not found.
//
// - `204 No Content` for success.
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
	err = c.apiModel.TransferApp(o)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}
