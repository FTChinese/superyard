package controller

import (
	"net/http"
	"strings"

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

// Exists tests if an account with the specified username or email exists
func (s StaffController) Exists(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))
	}

	key := req.Form.Get("k")
	val := req.Form.Get("v")

	if key == "" || val == "" {
		resp := util.NewBadRequest("Both 'k' and 'v' should be present in query string")
		view.Render(w, resp)

		return
	}

	var exists bool

	switch key {
	case "name":
		exists, err = s.model.StaffNameExists(val)
	case "email":
		exists, err = s.model.StaffEmailExists(val)
	// 400 Bad Request
	// {message: "..."}
	default:
		resp := util.NewBadRequest("The value of 'k' must be one of 'name' or 'email'")
		view.Render(w, resp)
		return
	}

	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}
	// 404 Not Found
	if !exists {
		view.Render(w, util.NewNotFound())

		return
	}

	view.Render(w, util.NewNoContent())
}

// Auth handles authentication process
// Input {userName: string, password: string, userIp: string}
func (s StaffController) Auth(w http.ResponseWriter, req *http.Request) {
	var login staff.Login

	// { message: "Problems parsing JSON" }
	if err := util.Parse(req.Body, &login); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	login.Sanitize()

	account, err := s.model.Auth(login)

	// { message: "xxxxx" } if server errored
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
// Input {email: string}
func (s StaffController) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	email, err := util.GetJSONString(req.Body, "email")

	// { message: "Problems parsing JSON" }
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// { message: "Validation failed"
	// 	 error: {
	//	    field: "email",
	//		code: "missing_field" | "invalid"
	//	 }
	// }
	if result := util.ValidateEmail(email); result.IsInvalid {
		view.Render(w, util.NewUnprocessable(result))

		return
	}

	err = s.model.RequestResetToken(email)

	// { message: "xxxxxxx" }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewNoContent())
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
func (s StaffController) VerifyToken(w http.ResponseWriter, req *http.Request) {
	token := chi.URLParam(req, "token")

	// { message: "Invalid request URI" }
	if token == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	account, err := s.model.VerifyResetToken(token)

	// 404 Not Found if the token does not exist
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// 200 OK { email: "foo@bar.org"}
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

	// { message: "Problems parsing JSON" }
	if err := util.Parse(req.Body, &reset); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	reset.Sanitize()

	// { message: "Validation failed" | "The length of password should not exceed 128 chars",
	// 	field: "password",
	//	code: "missing_field" | "invalid"
	// }
	if r := util.ValidatePassword(reset.Password); r.IsInvalid {
		resp := util.NewUnprocessable(r)
		view.Render(w, resp)

		return
	}

	err := s.model.ResetPassword(reset)

	// { message: "xxxxxxx" }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

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

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	resp := util.NewResponse().NoCache().SetBody(p)

	view.Render(w, resp)
}

// UpdateDisplayName lets user to change displayed name
// Input {displayName: string}, max 20 chars
func (s StaffController) UpdateDisplayName(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	displayName, err := util.GetJSONString(req.Body, "email")

	// { message: "Problems parsing JSON" }
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	displayName = strings.TrimSpace(displayName)

	// { message: "The length of displayName should not exceed 20 chars",
	// 	field: "displayName",
	//	code: "invalid"
	// }
	if r := util.ValidateMaxLen(displayName, 20, "displayName"); r.IsInvalid {
		resp := util.NewUnprocessable(r)

		view.Render(w, resp)

		return
	}

	if r := util.ValidateIsEmpty(displayName, "displayName"); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err = s.model.UpdateName(userName, displayName)

	// { message: "Validation failed",
	// 	field: "displayName",
	//	code: "already_exists"
	// }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "displayName"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// UpdateEmail lets user to change user name
// Input {email: string}, max 80 chars
func (s StaffController) UpdateEmail(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	email, err := util.GetJSONString(req.Body, "email")

	// { message: "Problems parsing JSON" }
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// { message: "Validation failed" | "The length of email should not exceed 20 chars"
	//	 field: "email",
	//	 code: "missing_field" | "invalid"
	// }
	if r := util.ValidateEmail(email); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err = s.model.UpdateEmail(userName, email)

	// { message: "Validation failed",
	// 	field: "email",
	//	code: "already_exists"
	// }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "email"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// UpdatePassword lets user to change user name
// Input {old: string, new: string}, max 128 chars
// The max length limit is random.
// Password actually should not have length limit.
// But hashing extremely long strings takes time.
func (s StaffController) UpdatePassword(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var p staff.Password

	// { message: "Problems parsing JSON" }
	if err := util.Parse(req.Body, &p); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	p.Sanitize()

	// { message: "Validation failed" | "Password should not execeed 128 chars"
	// 	 error: {
	//	    field: "password",
	//		code: "missing_field" | "invalid"
	//	 }
	// }
	if r := p.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := s.model.UpdatePassword(userName, p)

	// { message: "xxxxx" }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}

// ListMyft shows all ftc accounts associated with current user
func (s StaffController) ListMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myfts, err := s.model.ListMyft(userName)

	// Note there won't be SQLNoRows here since return data is an array.
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(myfts))
}

// AddMyft allows a logged in user to associate cms account with a ftc account
// Input {email: string, password} to verify that this user actually owns this ftc account
func (s StaffController) AddMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var c staff.MyftCredential

	// { message: "Problems parsing JSON" }
	if err := util.Parse(req.Body, &c); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	c.Sanitize()

	err := s.model.AddMyft(userName, c)

	// 404 Not Found if myft credentials are wrong
	// 422 if this ftc account might already exist:
	// { message: "Validation failed",
	// 	field: "email",
	//	code: "already_exists"
	// }
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "email"))

		return
	}

	view.Render(w, util.NewNoContent())
}

// RemoveMyft deletes a ftc account owned by current user
func (s StaffController) RemoveMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myftID := chi.URLParam(req, "id")

	// { message: "Invalid request URI" }
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := s.model.DeleteMyft(userName, myftID)

	// Any server error
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	view.Render(w, util.NewNoContent())
}
