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

// AdminRouter handle endpoints related to superuser administration
type AdminRouter struct {
	adminModel admin.Env
	staffModel staff.Env  // used by administrator to retrieve staff profile
	apiModel   ftcapi.Env // used to delete personal access tokens when removing a staff
}

// NewAdminRouter creates a new instance of AdminRouter
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
//		"email": "foo.bar@ftchinese.com",
//		"userName": "foo.bar",
//		"displayName": "Foo Bar",
//		"department": "tech",
//		"groupMembers": 3
//	}
//
// 400 Bad Request if request body cannot be parsed:
//	{
//		"message": "Problems parsing JSON"
//	}
func (m AdminRouter) NewStaff(w http.ResponseWriter, req *http.Request) {
	var a staff.Account

	if err := util.Parse(req.Body, &a); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	a.Sanitize()

	// 422 Unprocessable Entity:
	//	{
	//		message: "Validation failed" | "The length of email should not exceed 20 chars" | "The length of userName should be within 1 to 20 chars" | "The length of displayName should be within 1 to 20 chars"
	//		field: "email" | "userName" | "displayName",
	//		code: "missing_field" | "invalid"
	//	}
	if r := a.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := m.adminModel.NewStaff(a)

	// 422 Unprocessable Entity:
	//	{
	//		message: "Validation failed",
	// 		field: "email | userName | displayName",
	//		code: "already_exists"
	//	}
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 204 No Content if a new staff is created.
	view.Render(w, util.NewNoContent())
}

// StaffRoster lists all staff with pagination support.
//
//	GET `/admin/staff/roster?page=<number>`
//
// 400 Bad Request if query string cannot be parsed, query parameter `page` cannot be found, or is not a number.
//
// 200 OK with body:
//	[
//		{
//			"id": 1,
//			"email": "foo.bar@ftchinese.com",
//			"userName": "foo.bar",
//			"displayName": "Foo Bar",
//			"department": "tech",
//			"groupMembers": 3
//		}
//	]
func (m AdminRouter) StaffRoster(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request if query string cannot be parsed.
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))
		return
	}

	page, err := getQueryParam(req, "page").toInt()

	// 400 Bad Request if query parameter `page` cannot be found, or is not a number
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))
		return
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
//	GET `/admin/staff/profile/{name}`
//
// 400 Bad Request if url does not cotain `name` part.
// 	{
//		"message": "Invalid request URI"
//	}
//
// 404 Not Found if the requested user is not found
//
// 200 OK:
//	{
//		"id": "",
//		"userName": "",
// 		"email": "",
//		"isActive": "",
//		"displayName": "",
//		"department": "",
//		"groupMembers": "",
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

// ReinstateStaff restore a previously deleted staff.
//
//	PUT `/admin/staff/profile/{name}`
//
// 400 Bad Request:
// 	{
//		"message": "Invalid request URI"
//	}
//
// 204 No Content
func (m AdminRouter) ReinstateStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	// TODO: should check if the user acutually existed.
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
// Input:
//	{
//		"email": "required, max 20 chars",
//		"userName": "required, max 20 chars",
//		"displayName": "optional, max 20 chars",
//		"department": "optional",
//		"groupMembers": int
// 	}
//
// - 400 Bad Request if request URL does not contain `name` or request body cannot be parsed.
// 	{
//		"message": "Invalid request URI | Problems parsing JSON"
//	}
//
// - 422 Unprocessable Entity:
// if any of required fields is missing
//	{
//		"message": "Validation failed",
//		"field": "email | userName | displayName",
//		"code": "missing_field"
//	}
// if the fields are invalid:
//	{
//		"message": "The length of email should not exceed 20 chars | The length of userName should be within 1 to 20 chars | The length of displayName should be within 1 to 20 chars",
//		"field": "email | userName | displayName",
//		"code": "invalid"
//	}
// if `email`, `userName` or `displayName` is already taken:
//	{
//		"message": "Validation failed",
//		"field": "email | userName | displayName",
//		"code": "already_exists"
//	}
//
// - 204 No Content
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
// It also deletes all myft account associated with this staff;
// Unset vip of all related myft account;
// Remove all personal access token to access next-api;
// Remove all access token to access backyard-api
//
// 	DELETE `/admin/staff/profile/{name}?rmvip=true|false`
//
// - 400 Bad Request if request URL does not contain `name`, or query string `rmvip` exists but cannot be converted to bool.
// 	{
//		"message": "Invalid request URI"
//	}
//
// - 204 No Content
func (m AdminRouter) DeleteStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	rmVIP, err := getQueryParam(req, "rmvip").toBool()

	// 400 Bad Request
	if err != nil {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
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

// VIPRoster lists all ftc account with vip set to true.
//
//	GET `/admin/vip`
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
//	PUT `/admin/vip/{myftId}`
func (m AdminRouter) GrantVIP(w http.ResponseWriter, req *http.Request) {
	myftID := chi.URLParam(req, "id")

	// { message: "Invalid request URI" }
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.GrantVIP(myftID)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}

// RevokeVIP removes a ftc account from vip.
//
//	DELETE `/admin/vip/{myftId}`
func (m AdminRouter) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	myftID := chi.URLParam(req, "id")

	// { message: "Invalid request URI" }
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := m.adminModel.RevokeVIP(myftID)

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}
