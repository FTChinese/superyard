package controller

import (
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

const userNameKey = "X-User-Name"

// StaffController handles staff related actions like authentication, password reset, personal settings.
type StaffController struct {
	model staff.Env
}

// Auth handles authentication process
// Input {userName: string, password: string, userIp: string}
func (s StaffController) Auth(w http.ResponseWriter, req *http.Request) {
	var login staff.Login

	if err := util.Parse(req.Body, &login); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	login.Sanitize()

	account, err := s.model.Auth(login)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
// Input {email: string}
func (s StaffController) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	email, err := util.GetJSONString(req.Body, "email")

	if err != nil {
		view.Render(w, util.NewInvalidJSON(err))

		return
	}

	if err := util.ValidateEmail(email); err != nil {
		view.Render(w, util.NewUnprocessable("", err))

		return
	}

	err = s.model.RequestResetToken(email)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	view.Render(w, util.NewNoContent())
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
func (s StaffController) VerifyToken(w http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "token")

	if token == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	account, err := s.model.VerifyResetToken(token)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	resp := util.NewResponse().
		NoCache().
		SetBody(map[string]string{
			"email": account.Email,
		})
	view.Render(w, resp)
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
// Input {token: string, password: string}
func (s StaffController) ResetPassword(w http.ResponseWriter, req *http.Request) {
	var reset staff.PasswordReset

	if err := util.Parse(req.Body, &reset); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	reset.Sanitize()

	err := s.model.ResetPassword(reset)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	view.Render(w, util.NewNoContent())
}

// Profile shows a user's profile.
// Request header must contain `X-User-Name`
// There should be a middleware to check if `X-User-Name` exists
func (s StaffController) Profile(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	p, err := s.model.Profile(userName)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))

		return
	}

	resp := util.NewResponse().NoCache().SetBody(p)

	view.Render(w, resp)
}

// UpdateDisplayName lets user to change displayed name
// Input {displayName: string}
func (s StaffController) UpdateDisplayName(w http.ResponseWriter, req *http.Request) {
	// userName := req.Header.Get(userNameKey)
}

// UpdateEmail lets user to change user name
// Input {email: string}
func (s StaffController) UpdateEmail(w http.ResponseWriter, req *http.Request) {
	// userName := req.Header.Get(userNameKey)
}

// UpdatePassword lets user to change user name
// Input {old: string, new: string}
func (s StaffController) UpdatePassword(w http.ResponseWriter, req *http.Request) {
	// userName := req.Header.Get(userNameKey)
}
