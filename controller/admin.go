package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"gitlab.com/ftchinese/backyard-api/model"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/staff"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/util"
)

// AdminRouter responds to administration tasks performed by a superuser.
type AdminRouter struct {
	model   model.StaffEnv // used by administrator to retrieve staff profile
	search  model.SearchEnv
	postman postoffice.Postman
}

// NewAdminRouter creates a new instance of AdminRouter.
func NewAdminRouter(db *sql.DB, p postoffice.Postman) AdminRouter {
	return AdminRouter{
		model:   model.StaffEnv{DB: db},
		search:  model.SearchEnv{DB: db},
		postman: p,
	}
}

// Exists tests if an account with the specified userName or email exists
//
//	GET admin/account/exists?k={name|email}&v={value}
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
		exists, err = router.model.NameExists(val)
	case "email":
		exists, err = router.model.EmailExists(val)

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

// CreateAccount create a new account for a staff.
//
// 	POST /admin/accounts
func (router AdminRouter) CreateAccount(w http.ResponseWriter, req *http.Request) {
	var account staff.Account

	if err := gorest.ParseJSON(req.Body, &account); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
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

	err = router.model.CreateAccount(account, password)

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

// ListAccounts lists all staff. Pagination is supported.
//
//	GET /admin/accounts?page=<number>
func (router AdminRouter) ListStaff(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request if query string cannot be parsed.
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	page, _ := GetQueryParam(req, "page").ToInt()
	pagination := util.NewPagination(page, 20)

	accounts, err := router.model.ListAccounts(pagination)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().NoCache().SetBody(accounts))
}

// StaffProfile gets a staff's profile.
//
//	GET /admin/accounts/{name}
func (router AdminRouter) StaffProfile(w http.ResponseWriter, req *http.Request) {
	userName, err := GetURLParam(req, "name").ToString()

	// 400 Bad Request if url does not contain `name` part.
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p, err := router.model.Profile(userName)

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
//	PUT /admin/accounts/{name}
func (router AdminRouter) ReinstateStaff(w http.ResponseWriter, req *http.Request) {
	userName, err := GetURLParam(req, "name").ToString()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	err = router.model.ActivateStaff(userName)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}

// UpdateAccount updates a staff's account.
//
//	PATCH /admin/accounts/{name}
//
// Input {userName: string, email: string, displayName: string, department: string, groupMembers: number}
func (router AdminRouter) UpdateAccount(w http.ResponseWriter, req *http.Request) {
	userName, err := GetURLParam(req, "name").ToString()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var account staff.Account

	// 400 Bad Request
	if err := gorest.ParseJSON(req.Body, &account); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	account.Sanitize()

	// 422 Unprocessable Entity
	if r := account.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	err = router.model.UpdateAccount(userName, account)

	// 422 Unprocessable Entity: already_exists
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}

// DeleteStaff flags a staff as inactive.
//
// 	DELETE /admin/accounts/{name}
//
// Input {revokeVip: true | false}
func (router AdminRouter) DeleteStaff(w http.ResponseWriter, req *http.Request) {
	userName, err := GetURLParam(req, "name").ToString()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	result, err := GetJSONResult(req.Body, "revokeVip")

	var revokeVIP bool
	// rmVIP defaults to true.
	if err != nil || !result.Exists() {
		revokeVIP = true
	} else {
		revokeVIP = result.Bool()
	}

	// Removes a staff and optionally VIP status associated with this staff.
	err = router.model.RemoveStaff(userName, revokeVIP)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}

// ListVIP lists all ftc account granted vip.
//
//	GET /admin/vip
func (router AdminRouter) ListVIP(w http.ResponseWriter, req *http.Request) {
	myfts, err := router.model.ListVIP()

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(myfts))
}

// GrantVIP grants vip to an ftc account.
//
//	PUT /admin/vip/{email}
func (router AdminRouter) GrantVIP(w http.ResponseWriter, req *http.Request) {
	email, err := GetURLParam(req, "id").ToString()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	u, err := router.search.FindUserByEmail(email)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	err = router.model.GrantVIP(u.UserID)
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
//	DELETE /admin/vip/{email}
func (router AdminRouter) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	email, err := GetURLParam(req, "email").ToString()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	u, err := router.search.FindUserByEmail(email)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	err = router.model.RevokeVIP(u.UserID)
	// 500 Internal Server Error
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}
