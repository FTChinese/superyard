package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"gitlab.com/ftchinese/backyard-api/model"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/ftchinese/backyard-api/ftcapi"
	"gitlab.com/ftchinese/backyard-api/staff"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/util"
)

// AdminRouter responds to administration tasks performed by a superuser.
type AdminRouter struct {
	staffModel model.StaffEnv  // used by administrator to retrieve staff profile
	apiModel   ftcapi.Env // used to delete personal access tokens when removing a staff
	postman    postoffice.Postman
}

// NewAdminRouter creates a new instance of AdminRouter.
func NewAdminRouter(db *sql.DB, p postoffice.Postman) AdminRouter {
	staff := model.StaffEnv{DB: db}
	api := ftcapi.Env{DB: db}

	return AdminRouter{
		staffModel: staff,
		apiModel:   api,
		postman:    p,
	}
}

// Exists tests if an account with the specified userName or email exists
//
//	GET admin/staff/exists?k={name|email}&v={:value}
func (router AdminRouter) Exists(w http.ResponseWriter, req *http.Request) {

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
		exists, err = router.staffModel.StaffNameExists(val)
	case "email":
		exists, err = router.staffModel.StaffEmailExists(val)

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
func (router AdminRouter) SignUp(w http.ResponseWriter, req *http.Request) {
	var account staff.Account

	if err := util.Parse(req.Body, &account); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	account.Sanitize()

	// 422 Unprocessable Entity:
	if r := account.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))

		return
	}

	password, err := gorest.RandomHex(4)
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	err = router.staffModel.CreateAccount(account, password)

	// 422 Unprocessable Entity:
	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Field = "email"
			reason.Code = view.CodeAlreadyExists
			view.Render(w, view.NewUnprocessable(reason))

			return
		}

		view.Render(w, view.NewDBFailure(err))

		return
	}

	parcel, err := account.SignupParcel(password)
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	go router.postman.Deliver(parcel)

	// 204 No Content if account new staff is created.
	view.Render(w, view.NewNoContent())
}

// StaffRoster lists all staff. Pagination is supported.
//
//	GET /admin/staff/roster?page=<number>
func (router AdminRouter) StaffRoster(w http.ResponseWriter, req *http.Request) {
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

	accounts, err := router.staffModel.StaffRoster(page, 20)

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
func (router AdminRouter) StaffProfile(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request if url does not cotain `name` part.
	if userName == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	p, err := router.staffModel.Profile(userName)

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
func (router AdminRouter) ReinstateStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	err := router.staffModel.ActivateStaff(userName)

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
func (router AdminRouter) UpdateStaff(w http.ResponseWriter, req *http.Request) {
	userName := getURLParam(req, "name").toString()

	// 400 Bad Request
	if userName == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	var account staff.Account

	// 400 Bad Request
	if err := util.Parse(req.Body, &account); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	account.Sanitize()

	// 422 Unprocessable Entity
	if r := account.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))

		return
	}

	err := router.staffModel.UpdateStaff(userName, account)

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
func (router AdminRouter) DeleteStaff(w http.ResponseWriter, req *http.Request) {
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
	err = router.staffModel.RemoveStaff(userName, rmVIP)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// Removes any personal access token used for next-api created by this staff
	err = router.apiModel.RemovePersonalAccess(0, userName)

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
func (router AdminRouter) VIPRoster(w http.ResponseWriter, req *http.Request) {
	myfts, err := router.staffModel.VIPRoster()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(myfts))
}

// GrantVIP grants vip to an ftc account.
//
//	PUT /admin/vip/{myftId}
func (router AdminRouter) GrantVIP(w http.ResponseWriter, req *http.Request) {
	myftID := getURLParam(req, "id").toString()

	// 400 Bad Request
	if myftID == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	err := router.staffModel.GrantVIP(myftID)

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
func (router AdminRouter) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	myftID := chi.URLParam(req, "id")

	// 400 Bad Request
	if myftID == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	err := router.staffModel.RevokeVIP(myftID)

	// 500 Internal Server Error
	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}
