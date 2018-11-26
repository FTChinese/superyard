package controller

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/go-mail/mail"
	"gitlab.com/ftchinese/backyard-api/postman"
	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// StaffRouter responds to CMS login and personal settings.
type StaffRouter struct {
	model   staff.Env
	postman postman.Env
}

// NewStaffRouter creates a new instance of StaffController
func NewStaffRouter(db *sql.DB, dialer *mail.Dialer) StaffRouter {
	return StaffRouter{
		model:   staff.Env{DB: db},
		postman: postman.Env{Dialer: dialer},
	}
}

// Auth respond to login request.
//
// 	POST /staff/auth
func (r StaffRouter) Auth(w http.ResponseWriter, req *http.Request) {
	var login staff.Login

	// `400 Bad Request` if body content cannot be parsed as JSON
	if err := util.Parse(req.Body, &login); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	login.Sanitize()

	account, err := r.model.Auth(login)

	// `404 Not Found`
	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `200 OK`
	view.Render(w, util.NewResponse().NoCache().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST /staff/password-reset/letter
func (r StaffRouter) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	email, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request`
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// `422 Unprocessable Entity`
	if reason := util.RequireEmail(email); reason != nil {
		view.Render(w, util.NewUnprocessable(reason))

		return
	}

	parcel, err := r.model.RequestResetToken(email)

	// `404 Not Found`
	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	go r.postman.SendPasswordReset(parcel)

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
//
// 	GET /staff/password-reset/tokens/{token}
func (r StaffRouter) VerifyToken(w http.ResponseWriter, req *http.Request) {
	token := getURLParam(req, "token").toString()

	// `400 Bad Request`
	if token == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	account, err := r.model.VerifyResetToken(token)

	// `404 Not Found`
	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `200 OK`
	resp := util.NewResponse().
		NoCache().
		SetBody(map[string]string{
			"email": account.Email,
		})
	view.Render(w, resp)
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
//
//	POST /staff/password-reset
func (r StaffRouter) ResetPassword(w http.ResponseWriter, req *http.Request) {
	var reset staff.PasswordReset

	// `400 Bad Request`
	if err := util.Parse(req.Body, &reset); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	reset.Sanitize()

	// `422 Unprocessable Entity`
	if r := util.RequirePassword(reset.Password); r != nil {
		resp := util.NewUnprocessable(r)
		view.Render(w, resp)

		return
	}

	err := r.model.ResetPassword(reset)

	// 404 Not Found if the token is is expired or not found.
	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// Profile shows a user's profile.
// Request header must contain `X-User-Name`.
//
//	 GET /user/profile
func (r StaffRouter) Profile(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	p, err := r.model.Profile(userName)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `200 OK`
	resp := util.NewResponse().NoCache().SetBody(p)

	view.Render(w, resp)
}

// UpdateDisplayName lets user to change displayed name.
//
//	PATCH /user/display-name
func (r StaffRouter) UpdateDisplayName(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	displayName, err := util.GetJSONString(req.Body, "displayName")

	// `400 Bad Request`
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	displayName = strings.TrimSpace(displayName)

	// `422 Unprocessable Entity`
	if r := util.RequireNotEmptyWithMax(displayName, 255, "displayName"); r != nil {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err = r.model.UpdateName(userName, displayName)

	// `422 Unprocessable Entity` if this `displayName` already exists
	if err != nil {
		if util.IsAlreadyExists(err) {
			reason := util.NewReasonAlreadyExists("displayName")
			view.Render(w, util.NewUnprocessable(reason))
			return
		}
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// UpdateEmail lets user to change email.
//
//	PATCH /user/email
func (r StaffRouter) UpdateEmail(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	email, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request`
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// `422 Unprocessable Entity`
	if r := util.RequireEmail(email); r != nil {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err = r.model.UpdateEmail(userName, email)

	// `422 Unprocessable Entity`
	if err != nil {
		if util.IsAlreadyExists(err) {
			reason := util.NewReasonAlreadyExists("email")
			view.Render(w, util.NewUnprocessable(reason))
			return
		}
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `204 No Content` if updated successfully.
	view.Render(w, util.NewNoContent())
}

// UpdatePassword lets user to change password.
//
//	PATCH /user/password
func (r StaffRouter) UpdatePassword(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var p staff.Password

	// `400 Bad Request`
	if err := util.Parse(req.Body, &p); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	p.Sanitize()

	// `422 Unprocessable Entity`
	if r := p.Validate(); r != nil {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := r.model.UpdatePassword(userName, p)

	// `403 Forbidden` if old password is wrong
	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// ListMyft shows all ftc accounts associated with current user.
//
//	GET /user/myft
func (r StaffRouter) ListMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myfts, err := r.model.ListMyft(userName)

	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	// `200 OK`
	view.Render(w, util.NewResponse().NoCache().SetBody(myfts))
}

// AddMyft allows a logged in user to associate cms account with a ftc account.
//
//	POST /user/myft
func (r StaffRouter) AddMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var credential staff.MyftCredential

	// `400 Bad Request` for invalid JSON
	if err := util.Parse(req.Body, &credential); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	credential.Sanitize()

	err := r.model.AddMyft(userName, credential)

	// `404 Not Found` if `email` + `password` verification failed.
	// `422 Unprocessable Entity` if this ftc account already exist.
	if err != nil {
		if util.IsAlreadyExists(err) {
			reason := util.NewReasonAlreadyExists("email")
			view.Render(w, util.NewUnprocessable(reason))
			return
		}
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// DeleteMyft deletes a ftc account owned by current user.
//
//	DELETE /user/myft/{id}
func (r StaffRouter) DeleteMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myftID := getURLParam(req, "id").toString()

	// `400 Bad Request`
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := r.model.DeleteMyft(userName, myftID)

	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}
