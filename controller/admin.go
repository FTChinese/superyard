package controller

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/postman"

	"github.com/go-chi/chi"
	"github.com/go-mail/mail"
	"gitlab.com/ftchinese/backyard-api/ftcapi"
	"gitlab.com/ftchinese/backyard-api/staff"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/admin"
	"gitlab.com/ftchinese/backyard-api/util"
)

// AdminRouter responds to administration tasks performed by a superuser.
type AdminRouter struct {
	adminModel admin.Env
	staffModel staff.Env  // used by administrator to retrieve staff profile
	apiModel   ftcapi.Env // used to delete personal access tokens when removing a staff
	postman    postman.Env
}

// NewAdminRouter creates a new instance of AdminRouter.
func NewAdminRouter(db *sql.DB, dialer *mail.Dialer) AdminRouter {
	admin := admin.Env{DB: db}
	staff := staff.Env{DB: db}
	api := ftcapi.Env{DB: db}
	mailer := postman.Env{Dialer: dialer}

	return AdminRouter{
		adminModel: admin,
		staffModel: staff,
		apiModel:   api,
		postman:    mailer,
	}
}

// Exists tests if an account with the specified userName or email exists
//
//	GET admin/staff/exists?k={name|email}&v={:value}
func (r AdminRouter) Exists(w http.ResponseWriter, req *http.Request) {

	key := req.FormValue("k")
	val := req.FormValue("v")

	// `400 Bad Request`
	if key == "" || val == "" {
		resp := view.NewBadRequest("Both 'k' and 'v' should be present in query string")
		view.Render(w, resp)

		return
	}

	var exists bool
	var err error

	switch key {
	case "name":
		exists, err = r.staffModel.StaffNameExists(val)
	case "email":
		exists, err = r.staffModel.StaffEmailExists(val)

	// `400 Bad Request`
	default:
		resp := view.NewBadRequest("The value of 'k' must be one of 'name' or 'email'")
		view.Render(w, resp)
		return
	}

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}
	// `404 Not Found`
	if !exists {
		view.Render(w, view.NewNotFound())

		return
	}

	// `204 No Content` if the user exists.
	view.Render(w, view.NewNoContent())
}

// NewStaff create a new account for a staff.
//
// 	POST /admin/staff/new
func (r AdminRouter) NewStaff(w http.ResponseWriter, req *http.Request) {
	var a staff.Account

	if err := util.Parse(req.Body, &a); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	a.Sanitize()

	// 422 Unprocessable Entity:
	if r := a.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))

		return
	}

	parcel, err := r.adminModel.NewStaff(a)

	if util.IsAlreadyExists(err) {
		reason := view.NewReason()
		reason.Field = "email"
		reason.Code = view.CodeAlreadyExists
		view.Render(w, view.NewUnprocessable(reason))

		return
	}
	// 422 Unprocessable Entity:
	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	go r.postman.SendAccount(parcel)

	// 204 No Content if a new staff is created.
	view.Render(w, view.NewNoContent())
}

// StaffRoster lists all staff. Pagination is supported.
//
//	GET /admin/staff/roster?page=<number>
func (r AdminRouter) StaffRoster(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request if query string cannot be parsed.
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	page, err := getQueryParam(req, "page").toInt()

	if err != nil {
		page = 1
	}

	accounts, err := r.adminModel.StaffRoster(page, 20)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().NoCache().SetBody(accounts))
}

// StaffProfile gets a staff's profile.
//
//	GET /admin/staff/profile/{name}
func (r AdminRouter) StaffProfile(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request if url does not cotain `name` part.
	if userName == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	p, err := r.staffModel.Profile(userName)

	// 404 Not Found if the requested user is not found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().NoCache().SetBody(p))
}

// ReinstateStaff restore a previously deactivated staff.
//
//	PUT /admin/staff/profile/{name}
func (r AdminRouter) ReinstateStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	err := r.adminModel.ActivateStaff(userName)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}

// UpdateStaff updates a staff's profile.
//
//	PATCH /admin/staff/profile/{name}
func (r AdminRouter) UpdateStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	var a staff.Account

	// 400 Bad Request
	if err := util.Parse(req.Body, &a); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	a.Sanitize()

	// 422 Unprocessable Entity
	if r := a.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))

		return
	}

	err := r.adminModel.UpdateStaff(userName, a)

	// 422 Unprocessable Entity: already_exists
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
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
func (r AdminRouter) DeleteStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

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
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// Removes any personal access token used for next-api created by this staff
	err = r.apiModel.RemovePersonalAccess(0, userName)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}

// VIPRoster lists all ftc account granted vip.
//
//	GET /admin/vip
func (r AdminRouter) VIPRoster(w http.ResponseWriter, req *http.Request) {
	myfts, err := r.adminModel.VIPRoster()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(myfts))
}

// GrantVIP grants vip to an ftc account.
//
//	PUT /admin/vip/{myftId}
func (r AdminRouter) GrantVIP(w http.ResponseWriter, req *http.Request) {
	myftID := getURLParam(req, "id").toString()

	// 400 Bad Request
	if myftID == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	err := r.adminModel.GrantVIP(myftID)

	// 500 Internal Server Error
	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}

// RevokeVIP removes a ftc account from vip.
//
//	DELETE /admin/vip/{myftId}
func (r AdminRouter) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	myftID := chi.URLParam(req, "id")

	// 400 Bad Request
	if myftID == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	err := r.adminModel.RevokeVIP(myftID)

	// 500 Internal Server Error
	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}
