package controller

import (
	"database/sql"
	"net/http"
	"strings"

	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// StaffRouter handles staff related actions like authentication, password reset, personal settings.
type StaffRouter struct {
	model staff.Env
}

// NewStaffRouter creates a new instance of StaffController
func NewStaffRouter(db *sql.DB) StaffRouter {
	model := staff.Env{DB: db}

	return StaffRouter{
		model: model,
	}
}

// Exists tests if an account with the specified username or email exists
func (r StaffRouter) Exists(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
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
		exists, err = r.model.StaffNameExists(val)
	case "email":
		exists, err = r.model.StaffEmailExists(val)
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
// POST `/staff/auth`
// Input {userName: string, password: string, userIp: string}
func (r StaffRouter) Auth(w http.ResponseWriter, req *http.Request) {
	var login staff.Login

	// `400 Bad Request` if body content cannot be parsed as JSON
	// ```json
	// { "message": "Problems parsing JSON" }
	// ```
	if err := util.Parse(req.Body, &login); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	login.Sanitize()

	account, err := r.model.Auth(login)

	// `404 Not Found` if `userName` does not exist or `password` is wrong.
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `200 OK` with body:
	// ```json
	// {
	//	"id": 1,
	//	"email": "foo.bar@ftchinese.com",
	//	"userName": "foo.bar",
	//	"displayName": "Foo Bar",
	//	"department": "tech",
	//	"groupMembers": 3
	// }
	// ```
	view.Render(w, util.NewResponse().NoCache().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
// POST `/staff/password-reset/letter`
// Input `{ "email": "foo.bar@ftchinese.com" }`
func (r StaffRouter) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	email, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request` if request body cannot be parsed as JSON.
	// `{ "message": "Problems parsing JSON" }`
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// `422 Unprocessable Entity` if `email` is missing or invalid.
	// ```json
	// {
	//	"message": "Validation failed"
	//	"field": "email",
	// 	"code": "missing_field" | "invalid"
	// }
	// ```
	if result := util.ValidateEmail(email); result.IsInvalid {
		view.Render(w, util.NewUnprocessable(result))

		return
	}

	err = r.model.RequestResetToken(email)

	// `404 Not Found` if the `email` is not found
	//
	// `500 Internal Server Error` if token cannot be generated, or token cannot be saved, or email cannot be sent.
	// `{ "message": "xxxxxxx" }`
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// `204 No Content` if password reset letter is sent.
	view.Render(w, util.NewNoContent())
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
// GET `/staff/password-reset/tokens/{token}`
func (r StaffRouter) VerifyToken(w http.ResponseWriter, req *http.Request) {
	token := getURLParam(req, "token").toString()

	// `400 Bad Request` if request URL does not contain `{token}` part
	// `{ "message": "Invalid request URI" }`
	if token == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	account, err := r.model.VerifyResetToken(token)

	// `404 Not Found` if the token does not exist
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `200 OK` with body `{ "email": "foo@bar.org" }`
	resp := util.NewResponse().
		NoCache().
		SetBody(map[string]string{
			"email": account.Email,
		})
	view.Render(w, resp)
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
// POST `/staff/password-reset`
// Input `{ "token": string, "password": string }`
func (r StaffRouter) ResetPassword(w http.ResponseWriter, req *http.Request) {
	var reset staff.PasswordReset

	// `400 Bad Request` is request body cannot be parsed as JSON.
	// { "message": "Problems parsing JSON" }
	if err := util.Parse(req.Body, &reset); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	reset.Sanitize()

	// `422 Unprocessable Entity` if validation failed.
	// ```json
	// {
	//	"message": "Validation failed" | "The length of password should not exceed 128 chars",
	// 	"field": "password",
	//	"code": "missing_field" | "invalid"
	// }
	// ```
	if r := util.ValidatePassword(reset.Password); r.IsInvalid {
		resp := util.NewUnprocessable(r)
		view.Render(w, resp)

		return
	}

	err := r.model.ResetPassword(reset)

	// `500 Internal Server Error` if database errored.
	// `{ "message": "xxxxxxx" }`
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `204 No Content` if password is reset succesfully.
	view.Render(w, util.NewNoContent())
}

// Profile shows a user's profile.
// Request header must contain `X-User-Name`
// GET `/user/profile`
func (r StaffRouter) Profile(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	p, err := r.model.Profile(userName)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `200 OK` with body:
	// ```json
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
	// ```
	resp := util.NewResponse().NoCache().SetBody(p)

	view.Render(w, resp)
}

// UpdateDisplayName lets user to change displayed name
// PATCH `/user/display-name`
// Input `{ "displayName": "max 20 chars" }`
func (r StaffRouter) UpdateDisplayName(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	displayName, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request` if request body cannot be parsed as JSON
	// `{ "message": "Problems parsing JSON" }`
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	displayName = strings.TrimSpace(displayName)

	// `422 Unprocessable Entity`
	// ```json
	// {
	//		"message": "Validation failed | The length of displayName should not exceed 20 chars",
	// 		"field": "displayName",
	//		"code": "missing_field | invalid"
	// }
	// ```
	if r := util.ValidateIsEmpty(displayName, "displayName"); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	if r := util.ValidateMaxLen(displayName, 20, "displayName"); r.IsInvalid {
		resp := util.NewUnprocessable(r)

		view.Render(w, resp)

		return
	}

	err = r.model.UpdateName(userName, displayName)

	// `422 Unprocessable Entity` if this `displayName` already exists
	// ```json
	// {
	//		"message": "Validation failed",
	// 		"field": "displayName",
	//		"code": "already_exists"
	// }
	// ```
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "displayName"))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// UpdateEmail lets user to change user name
// PATCH `/user/email`
// Input `{ "email": "max 80 chars" }`
func (r StaffRouter) UpdateEmail(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	email, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request` if request body cannot be pased as JSON.
	// `{ "message": "Problems parsing JSON" }`
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// `422 Unprocessable Entity` for validation failure.
	// ```json
	// {
	//		message: "Validation failed | max 80 chars"
	//	 	field: "email",
	//	 	code: "missing_field | invalid"
	// }
	// ```
	if r := util.ValidateEmail(email); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err = r.model.UpdateEmail(userName, email)

	// `422 Unprocessable Entity` if the email to use already exists
	// ```json
	// {
	//		"message": "Validation failed",
	// 		"field": "email",
	//		"code": "already_exists"
	// }
	// ```
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "email"))

		return
	}

	// `204 No Content` if updated successfully.
	view.Render(w, util.NewNoContent())
}

// UpdatePassword lets user to change user name
// PATCH `/user/password`
// Input `{ "old": "max 128 chars", "new": "max 128 chars" }`,
// The max length limit is random.
// Password actually should not have length limit.
// But hashing extremely long strings takes time.
func (r StaffRouter) UpdatePassword(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var p staff.Password

	// `400 Bad Request` if request body cannot be parsed.
	// `{ "message": "Problems parsing JSON" }`
	if err := util.Parse(req.Body, &p); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	p.Sanitize()

	// `422 Unprocessable Entity` if either `old` or `new` is missing in request body, or password is too long.
	// ```json
	// {
	//		"message": "Validation failed | Password should not execeed 128 chars",
	//	    "field": "password",
	//		"code": "missing_field | invalid"
	// }
	// ```
	if r := p.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := r.model.UpdatePassword(userName, p)

	// `403 Forbidden` if old password is wrong
	// `{ "message": "wrong password" }``
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `204 No Content` if password updated successfully.
	view.Render(w, util.NewNoContent())
}

// ListMyft shows all ftc accounts associated with current user
// GET `/user/myft`
func (r StaffRouter) ListMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myfts, err := r.model.ListMyft(userName)

	// `500 Internal Server Error` if any server errored.
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// `200 OK`
	// ```json
	//	[{
	//		"myftId": "",
	//		"myftEmail": "",
	//		"isVip": "boolean"
	//	}]
	// ````
	view.Render(w, util.NewResponse().NoCache().SetBody(myfts))
}

// AddMyft allows a logged in user to associate cms account with a ftc account
// POST `/user/myft`
// Input `{ "email": "string", "password": "string" }`
func (r StaffRouter) AddMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	var credential staff.MyftCredential

	// `400 Bad Request` for invalid JSON
	// `{ "message": "Problems parsing JSON" }`
	if err := util.Parse(req.Body, &credential); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	credential.Sanitize()

	err := r.model.AddMyft(userName, credential)

	// `404 Not Found` if `email` + `password` verification failed.
	// `422 Unprocessable Entity` if this ftc account already exist.
	//	```json
	// {
	//		"message": "Validation failed",
	// 		"field": "email",
	//		"code": "already_exists"
	// }
	// ```
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "email"))

		return
	}

	// `204 No Content` if this ftc account is verified and associated with current user.
	view.Render(w, util.NewNoContent())
}

// DeleteMyft deletes a ftc account owned by current user
// DELETE `/user/myft/{id}`
func (r StaffRouter) DeleteMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myftID := getURLParam(req, "id").toString()

	// `400 Bad Request` if myft id is not present in URL.
	// `{ "message": "Invalid request URI" }``
	if myftID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	err := r.model.DeleteMyft(userName, myftID)

	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `204 No Content` if user removes this ftc account from hist myft account list.
	view.Render(w, util.NewNoContent())
}
