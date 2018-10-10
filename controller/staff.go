package controller

import (
	"database/sql"
	"net/http"
	"strings"

	"gitlab.com/ftchinese/backyard-api/staff"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// StaffRouter responds to CMS login and personal settings.
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

// Auth respond to login request.
//
// 	POST `/staff/auth`
//
// Input
// 	{
// 		"userName": "foo.bar",
// 		"password": "abcedfg",
// 		"userIp": "127.0.0.1"
// 	}
//
// - `400 Bad Request` if body content cannot be parsed as JSON
//	{
// 		"message": "Problems parsing JSON"
// 	}
//
// - `404 Not Found` if `userName` does not exist or `password` is wrong.
//
// - `200 OK` with body:
//	{
//		"id": 1,
//		"email": "foo.bar@ftchinese.com",
//		"userName": "foo.bar",
//		"displayName": "Foo Bar",
//		"department": "tech",
//		"groupMembers": 3
//	}
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
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `200 OK`
	view.Render(w, util.NewResponse().NoCache().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST `/staff/password-reset/letter`
//
// Input:
// 	{
// 		"email": "foo.bar@ftchinese.com"
// 	}
//
// - `400 Bad Request` if request body cannot be parsed as JSON.
//	{
//		"message": "Problems parsing JSON"
//	}
//
// - `422 Unprocessable Entity` if `email` is missing or invalid.
//	{
//		"message": "Validation failed"
//		"field": "email",
// 		"code": "missing_field" | "invalid"
//	}
//
// - `404 Not Found` if the `email` is not found.
//
// - `500 Internal Server Error` if token cannot be generated, or token cannot be saved, or email cannot be sent.
//	{
// 		"message": "xxxxxxx"
// 	}
//
// - `204 No Content` if password reset letter is sent.
func (r StaffRouter) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	email, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request`
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// `422 Unprocessable Entity`
	if result := util.ValidateEmail(email); result.IsInvalid {
		view.Render(w, util.NewUnprocessable(result))

		return
	}

	err = r.model.RequestResetToken(email)

	// `404 Not Found`
	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
//
// 	GET `/staff/password-reset/tokens/{token}`
//
// - `400 Bad Request` if request URL does not contain `token` part
//	{
//		"message": "Invalid request URI"
//	}
//
// - `404 Not Found` if the token does not exist
//
// - `200 OK` with body
// 	{
// 		"email": "foo.bar@ftchinese.com"
// 	}
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
		view.Render(w, util.NewDBFailure(err, ""))

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
//	POST `/staff/password-reset`
//
// Input:
// 	{
// 		"token": "reset token client extracted from url",
// 		"password": "8 to 128 chars"
// 	}
//
// - `400 Bad Request` if request body cannot be parsed as JSON.
//	{
// 		"message": "Problems parsing JSON"
//	}
//
// - `422 Unprocessable Entity` if validation failed.
// 	{
//		"message": "Validation failed | The length of password should not exceed 128 chars",
// 		"field": "password",
//		"code": "missing_field | invalid"
// 	}
//
// - 404 Not Found if the token is expired or not found.
//
// - `204 No Content` if password is reset succesfully.
func (r StaffRouter) ResetPassword(w http.ResponseWriter, req *http.Request) {
	var reset staff.PasswordReset

	// `400 Bad Request`
	if err := util.Parse(req.Body, &reset); err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	reset.Sanitize()

	// `422 Unprocessable Entity`
	if r := util.ValidatePassword(reset.Password); r.IsInvalid {
		resp := util.NewUnprocessable(r)
		view.Render(w, resp)

		return
	}

	err := r.model.ResetPassword(reset)

	// 404 Not Found if the token is is expired or not found.
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// Profile shows a user's profile.
// Request header must contain `X-User-Name`.
//
//	 GET `/user/profile`
//
// - `404 Not Found` if this user does not exist.
//
// - `200 OK` with body:
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
func (r StaffRouter) Profile(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	p, err := r.model.Profile(userName)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `200 OK`
	resp := util.NewResponse().NoCache().SetBody(p)

	view.Render(w, resp)
}

// UpdateDisplayName lets user to change displayed name.
//
//	PATCH `/user/display-name`
//
// Input:
// 	{
// 		"displayName": "max 20 chars"
// 	}
//
// - `400 Bad Request` if request body cannot be parsed as JSON.
//	{
// 		"message": "Problems parsing JSON"
//	}
//
// - `422 Unprocessable Entity` if validation failed:
// 	{
//		"message": "Validation failed | The length of displayName should not exceed 20 chars",
// 		"field": "displayName",
//		"code": "missing_field | invalid"
//	}
// if this `displayName` already exists
//	{
//		"message": "Validation failed",
// 		"field": "displayName",
//		"code": "already_exists"
//	}
//
// - `204 No Content` for success
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
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "displayName"))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// UpdateEmail lets user to change email.
//
//	PATCH `/user/email`
//
// Input
// 	{
// 		"email": "max 20 chars"
// 	}
//
// - `400 Bad Request` if request body cannot be parsed as JSON.
//	{
// 		"message": "Problems parsing JSON"
//	}
//
// - `422 Unprocessable Entity` for validation failure:
//	{
//		message: "Validation failed | The length of email should not exceed 20 chars"
//	 	field: "email",
//	 	code: "missing_field | invalid"
//	}
// if the email to use already exists
// {
//		"message": "Validation failed",
// 		"field": "email",
//		"code": "already_exists"
// }
//
// - `204 No Content` for success.
func (r StaffRouter) UpdateEmail(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	email, err := util.GetJSONString(req.Body, "email")

	// `400 Bad Request`
	if err != nil {
		view.Render(w, util.NewBadRequest(""))

		return
	}

	// `422 Unprocessable Entity`
	if r := util.ValidateEmail(email); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err = r.model.UpdateEmail(userName, email)

	// `422 Unprocessable Entity`
	if err != nil {
		view.Render(w, util.NewDBFailure(err, "email"))

		return
	}

	// `204 No Content` if updated successfully.
	view.Render(w, util.NewNoContent())
}

// UpdatePassword lets user to change password.
//
//	PATCH `/user/password`
//
// Input
// 	{
// 		"old": "max 128 chars",
// 		"new": "max 128 chars"
// 	}
// The max length limit is random.
// Password actually should not have length limit.
// But hashing extremely long strings takes time.
//
// - `400 Bad Request` if request body cannot be parsed as JSON.
//	{
// 		"message": "Problems parsing JSON"
//	}
//
// - `422 Unprocessable Entity` if either `old` or `new` is missing in request body, or password is too long.
//	{
//		"message": "Validation failed | Password should not execeed 128 chars",
//	    "field": "password",
//		"code": "missing_field | invalid"
//	}
//
// - `403 Forbidden` if old password is wrong
//	{
// 		"message": "wrong password"
// 	}
//
// - `204 No Content` for success.
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
	if r := p.Validate(); r.IsInvalid {
		view.Render(w, util.NewUnprocessable(r))

		return
	}

	err := r.model.UpdatePassword(userName, p)

	// `403 Forbidden` if old password is wrong
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// ListMyft shows all ftc accounts associated with current user.
//
//	GET `/user/myft`
//
// - `200 OK`
//	[
//		{
//			"myftId": "",
//			"myftEmail": "",
//			"isVip": "boolean"
//		}
//	]
func (r StaffRouter) ListMyft(w http.ResponseWriter, req *http.Request) {
	userName := req.Header.Get(userNameKey)

	myfts, err := r.model.ListMyft(userName)

	// `500 Internal Server Error`
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	// `200 OK`
	view.Render(w, util.NewResponse().NoCache().SetBody(myfts))
}

// AddMyft allows a logged in user to associate cms account with a ftc account.
//
//	POST `/user/myft`
//
// Input
// 	{
// 		"email": "string",
// 		"password": "string"
// 	}
//
// - `400 Bad Request` if request body cannot be parsed as JSON.
//	{
// 		"message": "Problems parsing JSON"
//	}
//
// - `404 Not Found` if `email` + `password` verification failed.
//
// - `422 Unprocessable Entity` if the ftc account to add already exist.
// {
//		"message": "Validation failed",
// 		"field": "email",
//		"code": "already_exists"
// }
//
// - `204 No Content`
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
		view.Render(w, util.NewDBFailure(err, "email"))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}

// DeleteMyft deletes a ftc account owned by current user.
//
//	DELETE `/user/myft/{id}`
//
// - `400 Bad Request` if request URL does not contain `id` part
//	{
//		"message": "Invalid request URI"
//	}
//
// - `204 No Content` for success
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
		view.Render(w, util.NewDBFailure(err, ""))

		return
	}

	// `204 No Content`
	view.Render(w, util.NewNoContent())
}
