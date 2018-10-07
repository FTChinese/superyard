package controller

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/ftchinese/backyard-api/ftcapi"
	"gitlab.com/ftchinese/backyard-api/staff"

	"gitlab.com/ftchinese/backyard-api/admin"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// AdminRouter responds to administration tasks performed by a superuser.
type AdminRouter struct {
	adminModel admin.Env
	staffModel staff.Env  // used by administrator to retrieve staff profile
	apiModel   ftcapi.Env // used to delete personal access tokens when removing a staff
}

// NewAdminRouter creates a new instance of AdminRouter.
func NewAdminRouter(db *sql.DB) AdminRouter {
	admin := admin.Env{DB: db}
	staff := staff.Env{DB: db}
	api := ftcapi.Env{DB: db}

	return AdminRouter{
		adminModel: admin,
		staffModel: staff,
		apiModel:   api,
	}
}

// Exists tests if an account with the specified userName or email exists
//
//	GET `/staff/exists?k={name|email}&v={:value}`
//
// - `400 Bad Request` if url query string cannot be parsed:
// 	{
// 		"message": "Bad request"
// 	}
// or either `k` or `v` cannot be found in query string:
// 	{
// 		"message": "Both 'k' and 'v' should be present in query string"
// 	}
// or if the value of url query parameter `k` is neither `name` nor `email`
// 	{
// 		"message": "The value of 'k' must be one of 'name' or 'email'"
// 	}
//
// - `404 Not Found` if the the user with the specified `name` or `email` is not found.
//
// - `204 No Content` if the user exists.
func (r AdminRouter) Exists(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	key := req.Form.Get("k")
	val := req.Form.Get("v")

	// `400 Bad Request`
	if key == "" || val == "" {
		resp := util.NewBadRequest("Both 'k' and 'v' should be present in query string")
		view.Render(w, resp)

		return
	}

	var exists bool

	switch key {
	case "name":
		exists, err = r.staffModel.StaffNameExists(val)
	case "email":
		exists, err = r.staffModel.StaffEmailExists(val)

	// `400 Bad Request`
	default:
		resp := util.NewBadRequest("The value of 'k' must be one of 'name' or 'email'")
		view.Render(w, resp)
		return
	}

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}
	// `404 Not Found`
	if !exists {
		view.Render(w, util.NewNotFound())

		return
	}

	// `204 No Content` if the user exists.
	view.Render(w, util.NewNoContent())
}

// NewStaff create a new account for a staff.
//
// 	POST /admin/staff/new
//
// Input:
//	{
//		"email": "foo.bar@ftchinese.com", // required, unique, max 80 chars
//		"userName": "foo.bar", // required, unique, max 255 chars
//		"displayName": "Foo Bar", // optional, unique, max 255 chars
//		"department": "tech", // optinal, max 255 chars
//		"groupMembers": 3  // required, > 0
//	}
//
// - 400 Bad Request if request body cannot be parsed:
//	{
//		"message": "Problems parsing JSON"
//	}
//
// - 422 Unprocessable Entity:
//
// if any of the required fields is missing
// 	{
// 		"message": "Validation failed",
// 		"field": "email | userName | groupMembers",
// 		"code": "missing"
// 	}
// if email is not a valid email address
// 	{
// 		"message": "Validation failed",
// 		"field": "email",
// 		"code": "invalid"
// 	}
// if the length of any string fields is over 255:
// 	{
// 		"message": "The length of xxx should not exceed 255 chars",
// 		"field": "email | userName | displayName | department",
// 		"code": "invalid"
// 	}
// if any of unique fields is already taken by others:
//	{
//		message: "Validation failed",
// 		field: "email | userName | displayName",
//		code: "already_exists"
//	}
//
// - 204 No Content if a new staff is created.
func (r AdminRouter) NewStaff(w http.ResponseWriter, req *http.Request) {
	var a staff.Account

	if err := util.Parse(req.Body, &a); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	a.Sanitize()

	// 422 Unprocessable Entity:
	if r := a.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := r.adminModel.NewStaff(a)

	// 422 Unprocessable Entity:
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 204 No Content if a new staff is created.
	view.Render(w, util.NewNoContent())
}

// StaffRoster lists all staff. Pagination is supported.
//
//	GET /admin/staff/roster?page=<number>
//
// `page` defaults to 1 if omitted or is not a number. Returns 20 entires per page.
//
// - 200 OK with an array:
//	[{
//		"id": 1,
//		"email": "foo.bar@ftchinese.com",
//		"userName": "foo.bar",
//		"displayName": "Foo Bar",
//		"department": "tech",
//		"groupMembers": 3
//	}]
func (r AdminRouter) StaffRoster(w http.ResponseWriter, req *http.Request) {
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

	accounts, err := r.adminModel.StaffRoster(page, 20)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 200 OK
	view.Render(w, util.NewResponse().NoCache().SetBody(accounts))
}

// StaffProfile gets a staff's profile.
//
//	GET /admin/staff/profile/{name}
//
// - 400 Bad Request if url does not contain the `name` part.
// 	{
//		"message": "Invalid request URI"
//	}
//
// - 404 Not Found if the requested user is not found
//
// - 200 OK:
//	{
//		"id": "",
//		"userName": "",
// 		"email": "",
//		"isActive": true,
//		"displayName": "",
//		"department": "",
//		"groupMembers": 3,
//		"createdAt": "",
//		"deactivatedAt": "",
//		"updatedAt": "",
//		"lastLoginAt": "",
//		"lastLoginIp": ""
//	}
func (r AdminRouter) StaffProfile(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request if url does not cotain `name` part.
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	p, err := r.staffModel.Profile(userName)

	// 404 Not Found if the requested user is not found
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// 200 OK
	view.Render(w, util.NewResponse().NoCache().SetBody(p))
}

// ReinstateStaff restore a previously deactivated staff.
//
//	PUT /admin/staff/profile/{name}
//
// - 400 Bad Request  if url does not contain the `name` part.
// 	{
//		"message": "Invalid request URI"
//	}
//
// - 204 No Content
func (r AdminRouter) ReinstateStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := r.adminModel.ActivateStaff(userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}

// UpdateStaff updates a staff's profile.
//
//	PATCH `/admin/staff/profile/{name}`
//
// Input and response are identical to creating a new staff `POST /admin/staff/new`.
func (r AdminRouter) UpdateStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	var a staff.Account

	// 400 Bad Request
	if err := util.Parse(req.Body, &a); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	a.Sanitize()

	// 422 Unprocessable Entity
	if r := a.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := r.adminModel.UpdateStaff(userName, a)

	// 422 Unprocessable Entity: already_exists
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}

// DeleteStaff flags a staff as inactive.
// It performs mutilple actions:
//
// 1. Turns the staff to inactive state so that he/she could no longer login to CMS;
//
// 2. Revoke VIP from all ftc account associated with this staff;
//
// 3. Unlink ftc accounts this staff previously linked with CMS account;
//
// 4. Remove all personal access tokens to access next-api;
//
// 5. Remove all access tokens to access backyard-api
//
// 	DELETE /admin/staff/profile/{name}?rmvip=<true|false>
// `rmvip` defaults to true if omitted, or cannot be converted to a boolean value.
//
// `name` is a staff's login name.
//
// - 400 Bad Request if request URL does not contain `name`.
// 	{
//		"message": "Invalid request URI"
//	}
//
// - 204 No Content for success.
func (r AdminRouter) DeleteStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	rmVIP, err := getQueryParam(req, "rmvip").toBool()

	// rmVIP defaults to true.
	if err != nil {
		rmVIP = true
	}

	// Removes a staff and optionally VIP status associated with this staff.
	err = r.adminModel.RemoveStaff(userName, rmVIP)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// Removes any personal access token used for next-api created by this staff
	err = r.apiModel.RemovePersonalAccess(0, userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}

// VIPRoster lists all ftc account granted vip.
//
//	GET /admin/vip
// - 200 OK with body:
//	[{
// 		"myftId": "string",
// 		"myftEmail": "string"
//	}]
func (r AdminRouter) VIPRoster(w http.ResponseWriter, req *http.Request) {
	myfts, err := r.adminModel.VIPRoster()

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(myfts))
}

// GrantVIP grants vip to a ftc account.
//
//	PUT /admin/vip/{myftId}
//
// - `400 Bad Request` if `myftId` is not present in URL.
//	{
//		"message": "Invalid request URI"
//	}
//
// - 204 No Content if granted.
func (r AdminRouter) GrantVIP(w http.ResponseWriter, req *http.Request) {
	myftID := getURLParam(req, "id").toString()

	// 400 Bad Request
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := r.adminModel.GrantVIP(myftID)

	// 500 Internal Server Error
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}

// RevokeVIP removes a ftc account from vip.
//
//	DELETE /admin/vip/{myftId}
//
// - `400 Bad Request` if `myftId` is not present in URL.
//	{
//		"message": "Invalid request URI"
//	}
//
// - 204 No Content if revoked successuflly.
func (r AdminRouter) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	myftID := chi.URLParam(req, "id")

	// 400 Bad Request
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := r.adminModel.RevokeVIP(myftID)

	// 500 Internal Server Error
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}
