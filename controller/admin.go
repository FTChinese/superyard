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

// AdminRouter handles endpoints under `/admin` path.
// All endpoints performs administration tasks:
//
// - POST `/admin/staff/new` creates a new staff;
//
// - GET `/admin/staff/roster?page=<number>` show the list of all staff;
//
// - GET `/admin/staff/profile/{name}` Show a staff's profile
//
// - PUT `/admin/staff/profile/{name}` Reinstate a previously deleted staff;
//
// - PATCH `/admin/staff/profile/{name}` Update staff's profile
//
// - DELETE `/admin/staff/profile/{name}?rmvip=true|false` Delete a staff
//
// - GET `/admin/vip` Show all myft accounts that are granted VIP.
//
// - PUT `/admin/vip/{myftId}` Grant vip to a myft account
//
// - DELETE `/admin/vip/{myftId}` Revoke vip status of a ftc account
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

// NewStaff create a new account for a staff.
//
// 	POST /admin/staff/new
//
// Input:
//	{
//		"email": "foo.bar@ftchinese.com", // required, max 80 chars, unique
//		"userName": "foo.bar", // required, max 20 chars, unique
//		"displayName": "Foo Bar", // optional, max 20 chars, unique
//		"department": "tech", // optinal, max 80 chars
//		"groupMembers": 3  // required
//	}
//
// - 400 Bad Request if request body cannot be parsed:
//	{
//		"message": "Problems parsing JSON"
//	}
//
// - 422 Unprocessable Entity:
//
// if email is missing
// 	{
// 		"message": "Validation failed",
// 		"field": "email",
// 		"code": "missing"
// 	}
// if email is not a valid email address
// 	{
// 		"message": "Validation failed",
// 		"field": "email",
// 		"code": "invalid"
// 	}
// if the length of email is over 80:
// 	{
// 		"message": "The length of email should not exceed 80 chars",
// 		"field": "email",
// 		"code": "invalid"
// 	}
// if userName is missing:
// 	{
// 		"message": "Validation failed",
// 		"field": "userName",
// 		"code": "missing"
// 	}
// if the length of userName is over 20:
// 	{
// 		"message": "The length of userName should not exceed 20 chars",
// 		"field": "userName",
// 		"code": "invalid"
// 	}
// if the length of displayName is over 20:
//	{
//		message: "The length of displayName should not exceed 20 chars"
//		field: "displayName",
//		code: "invalid"
//	}
// if any of email, userName or displayName is already taken by others:
//	{
//		message: "Validation failed",
// 		field: "email | userName | displayName",
//		code: "already_exists"
//	}
//
// - 204 No Content if a new staff is created.
func (m AdminRouter) NewStaff(w http.ResponseWriter, req *http.Request) {
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

	err := m.adminModel.NewStaff(a)

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
func (m AdminRouter) StaffRoster(w http.ResponseWriter, req *http.Request) {
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

	accounts, err := m.adminModel.StaffRoster(page, 20)

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
func (m AdminRouter) StaffProfile(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request if url does not cotain `name` part.
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	p, err := m.staffModel.Profile(userName)

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
func (m AdminRouter) ReinstateStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.ActivateStaff(userName)

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
func (m AdminRouter) UpdateStaff(w http.ResponseWriter, req *http.Request) {
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

	err := m.adminModel.UpdateStaff(userName, a)

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
func (m AdminRouter) DeleteStaff(w http.ResponseWriter, req *http.Request) {
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
	err = m.adminModel.RemoveStaff(userName, rmVIP)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// Removes any personal access token used for next-api created by this staff
	err = m.apiModel.RemovePersonalAccess(0, userName)

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
func (m AdminRouter) VIPRoster(w http.ResponseWriter, req *http.Request) {
	myfts, err := m.adminModel.VIPRoster()

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
func (m AdminRouter) GrantVIP(w http.ResponseWriter, req *http.Request) {
	myftID := getURLParam(req, "id").toString()

	// 400 Bad Request
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.GrantVIP(myftID)

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
func (m AdminRouter) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	myftID := chi.URLParam(req, "id")

	// 400 Bad Request
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.RevokeVIP(myftID)

	// 500 Internal Server Error
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 204 No Content
	view.Render(w, util.NewNoContent())
}
