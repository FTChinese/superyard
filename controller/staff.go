package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"gitlab.com/ftchinese/backyard-api/model"
	"gitlab.com/ftchinese/backyard-api/types/user"
	"net/http"
	"strings"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/types/util"
)

// StaffRouter responds to CMS login and personal settings.
type StaffRouter struct {
	model   model.StaffEnv
	search  model.SearchEnv
	postman postoffice.Postman
}

// NewStaffRouter creates a new instance of StaffController
func NewStaffRouter(db *sql.DB, p postoffice.Postman) StaffRouter {
	return StaffRouter{
		model:   model.StaffEnv{DB: db},
		search:  model.SearchEnv{DB: db},
		postman: p,
	}
}

// Login verifies a user's user name and password. Headers: `X-User-Ip`.
//
// 	POST /staff/login
//
// Input {userName: string, password: string}
// Response 204 No Content if password and user name combination matched.
// Client should then proceed to fetch this user's account data.
func (router StaffRouter) Login(w http.ResponseWriter, req *http.Request) {
	var login staff.Login

	// `400 Bad Request` if body content cannot be parsed as JSON
	if err := gorest.ParseJSON(req.Body, &login); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	login.Sanitize()

	matched, err := router.model.IsPasswordMatched(login.UserName, login.Password)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	if !matched {
		view.Render(w, view.NewForbidden("wrong credentials"))
		return
	}

	userIP := req.Header.Get("X-User-Ip")
	go router.model.UpdateLoginHistory(login, userIP)

	// `200 OK`
	view.Render(w, view.NewNoContent())
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST /staff/password-reset/letter
//
// Input {email: string}
func (router StaffRouter) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	result, err := GetJSONResult(req.Body, "email")
	// `400 Bad Request`
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	email := strings.TrimSpace(result.String())
	// `422 Unprocessable Entity`
	if reason := util.RequireEmail(email); reason != nil {
		view.Render(w, view.NewUnprocessable(reason))
		return
	}

	account, err := router.model.LoadAccountByEmail(email, true)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	th, err := account.TokenHolder()
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	err = router.model.SavePwResetToken(th)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	parcel, err := account.PasswordResetParcel(th.GetToken())
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	go router.postman.Deliver(parcel)

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
//
// 	GET /staff/password-reset/tokens/{token}
//
// Output {email: string}
func (router StaffRouter) VerifyToken(w http.ResponseWriter, req *http.Request) {
	token, err := GetURLParam(req, "token").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	account, err := router.model.VerifyResetToken(token)

	// `404 Not Found`
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `200 OK`
	resp := view.NewResponse().
		SetBody(map[string]string{
			"email": account.Email,
		})
	view.Render(w, resp)
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
//
//	POST /staff/password-reset
//
// Input {token: string, password: string}
func (router StaffRouter) ResetPassword(w http.ResponseWriter, req *http.Request) {
	var reset staff.PasswordReset

	// `400 Bad Request`
	if err := gorest.ParseJSON(req.Body, &reset); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	reset.Sanitize()

	// `422 Unprocessable Entity`
	if r := util.RequirePassword(reset.Password); r != nil {
		resp := view.NewUnprocessable(r)
		view.Render(w, resp)

		return
	}

	err := router.model.ResetPassword(reset)

	// 404 Not Found if the token is is expired or not found.
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}

// Account loads a staff's account data. Header `X-User-Name`
//
//	GET /staff/account
func (router StaffRouter) Account(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	a, err := router.model.LoadAccountByName(userName, true)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(a))
}

// Profile shows a staff's profile.
// Header `X-User-Name`.
//
//	 GET /staff/profile
func (router StaffRouter) Profile(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	p, err := router.model.Profile(userName)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(p))
}

// UpdateDisplayName lets user to change displayed name.
//
//	PATCH /staff/display-name
//
// Input {displayName: string}
func (router StaffRouter) UpdateDisplayName(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	result, err := GetJSONResult(req.Body, "displayName")
	// `400 Bad Request`
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	displayName := strings.TrimSpace(result.String())

	// `422 Unprocessable Entity`
	if r := util.RequireNotEmptyWithMax(displayName, 255, "displayName"); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	err = router.model.UpdateName(userName, displayName)

	// `422 Unprocessable Entity` if this `displayName` already exists
	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Code = view.CodeAlreadyExists
			reason.Field = "displayName"
			view.Render(w, view.NewUnprocessable(reason))
			return
		}
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}

// UpdateEmail lets user to change email.
//
//	PATCH /staff/email
//
// Input {email: string}
func (router StaffRouter) UpdateEmail(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	result, err := GetJSONResult(req.Body, "email")
	// `400 Bad Request`
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	email := strings.TrimSpace(result.String())
	// `422 Unprocessable Entity`
	if r := util.RequireEmail(email); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	err = router.model.UpdateEmail(userName, email)

	// `422 Unprocessable Entity`
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

	// `204 No Content` if updated successfully.
	view.Render(w, view.NewNoContent())
}

// UpdatePassword lets user to change password.
//
//	PATCH /staff/password
//
// Input {oldPassword: string, newPassword: string}
func (router StaffRouter) UpdatePassword(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var p staff.Password

	// `400 Bad Request`
	if err := gorest.ParseJSON(req.Body, &p); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p.Sanitize()

	// `422 Unprocessable Entity`
	if r := p.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	err := router.model.UpdatePassword(userName, p)

	// `403 Forbidden` if old password is wrong
	if err != nil {
		if err == util.ErrWrongPassword {
			view.Render(w, view.NewForbidden(err.Error()))
			return
		}
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}

// AddMyft allows a logged in user to associate cms account with a ftc account.
//
//	POST /staff/myft
//
// Input {email: string, password: string}
func (router StaffRouter) AddMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var login user.Login

	// `400 Bad Request` for invalid JSON
	if err := gorest.ParseJSON(req.Body, &login); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	login.Sanitize()

	err := router.model.AddMyft(userName, login)

	// `404 Not Found` if `email` + `password` verification failed.
	// `422 Unprocessable Entity` if this ftc account already exist.
	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Code = view.CodeAlreadyExists
			reason.Field = "email"
			view.Render(w, view.NewUnprocessable(reason))
			return
		}
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}

// ListMyft shows all ftc accounts associated with current user.
//
//	GET /staff/myft
func (router StaffRouter) ListMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myfts, err := router.model.ListMyft(userName)

	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `200 OK`
	view.Render(w, view.NewResponse().NoCache().SetBody(myfts))
}

// DeleteMyft deletes an ftc account owned by current user.
//
//	DELETE /staff/myft
//
// Input {email: string}
func (router StaffRouter) DeleteMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	result, err := GetJSONResult(req.Body, "email")
	// `400 Bad Request`
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	email := result.String()

	if r := util.RequireEmail(email); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	u, err := router.search.FindUserByEmail(email)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	err = router.model.DeleteMyft(userName, u.UserID)

	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}
