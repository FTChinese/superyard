package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest/postoffice"
	"gitlab.com/ftchinese/backyard-api/model"
	"net/http"
	"strings"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
)

// StaffRouter responds to CMS login and personal settings.
type StaffRouter struct {
	model   model.StaffEnv
	postman postoffice.Postman
}

// NewStaffRouter creates a new instance of StaffController
func NewStaffRouter(db *sql.DB, p postoffice.Postman) StaffRouter {
	return StaffRouter{
		model:   model.StaffEnv{DB: db},
		postman: p,
	}
}

// Auth respond to login request.
//
// 	POST /staff/auth
func (router StaffRouter) Auth(w http.ResponseWriter, req *http.Request) {
	var login staff.Login

	// `400 Bad Request` if body content cannot be parsed as JSON
	if err := util.Parse(req.Body, &login); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	login.Sanitize()

	account, err := router.model.Auth(login)

	// `404 Not Found`
	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// `200 OK`
	view.Render(w, view.NewResponse().NoCache().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST /staff/password-reset/letter
func (router StaffRouter) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	email, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request`
	if err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	// `422 Unprocessable Entity`
	if reason := util.RequireEmail(email); reason != nil {
		view.Render(w, view.NewUnprocessable(reason))

		return
	}

	acnt, err := router.model.FindAccountByEmail(email, true)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	th, err := acnt.TokenHolder()
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	err = router.model.SavePwResetToken(th)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	parcel, err := acnt.PasswordResetParcel(th.GetToken())
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
func (router StaffRouter) VerifyToken(w http.ResponseWriter, req *http.Request) {
	token := getURLParam(req, "token").toString()

	// `400 Bad Request`
	if token == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	acnt, err := router.model.VerifyResetToken(token)

	// `404 Not Found`
	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// `200 OK`
	resp := view.NewResponse().
		NoCache().
		SetBody(map[string]string{
			"email": acnt.Email,
		})
	view.Render(w, resp)
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
//
//	POST /staff/password-reset
func (router StaffRouter) ResetPassword(w http.ResponseWriter, req *http.Request) {
	var reset staff.PasswordReset

	// `400 Bad Request`
	if err := util.Parse(req.Body, &reset); err != nil {
		view.Render(w, view.NewBadRequest(""))

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

// Profile shows a user's profile.
// Request header must contain `X-User-Name`.
//
//	 GET /user/profile
func (router StaffRouter) Profile(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	p, err := router.model.Profile(userName)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// `200 OK`
	resp := view.NewResponse().NoCache().SetBody(p)

	view.Render(w, resp)
}

// UpdateDisplayName lets user to change displayed name.
//
//	PATCH /user/display-name
func (router StaffRouter) UpdateDisplayName(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	displayName, err := util.GetJSONString(req.Body, "displayName")

	// `400 Bad Request`
	if err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	displayName = strings.TrimSpace(displayName)

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
//	PATCH /user/email
func (router StaffRouter) UpdateEmail(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	email, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request`
	if err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

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
//	PATCH /user/password
func (router StaffRouter) UpdatePassword(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var p staff.Password

	// `400 Bad Request`
	if err := util.Parse(req.Body, &p); err != nil {
		view.Render(w, view.NewBadRequest(""))

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
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}

// ListMyft shows all ftc accounts associated with current user.
//
//	GET /user/myft
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

// AddMyft allows a logged in user to associate cms account with a ftc account.
//
//	POST /user/myft
func (router StaffRouter) AddMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var credential staff.MyftCredential

	// `400 Bad Request` for invalid JSON
	if err := util.Parse(req.Body, &credential); err != nil {
		view.Render(w, view.NewBadRequest(""))

		return
	}

	credential.Sanitize()

	err := router.model.AddMyft(userName, credential)

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

// DeleteMyft deletes a ftc account owned by current user.
//
//	DELETE /user/myft/{id}
func (router StaffRouter) DeleteMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myftID := getURLParam(req, "id").toString()

	// `400 Bad Request`
	if myftID == "" {
		view.Render(w, view.NewBadRequest("Invalid request URI"))

		return
	}

	err := router.model.DeleteMyft(userName, myftID)

	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, view.NewDBFailure(err))

		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}
