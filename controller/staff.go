package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/employee"
	"gitlab.com/ftchinese/backyard-api/models/reader"
	"gitlab.com/ftchinese/backyard-api/models/util"
	"gitlab.com/ftchinese/backyard-api/repository/staff"
	"net/http"
)

// StaffRouter responds to CMS login and personal settings.
type StaffRouter struct {
	postman postoffice.Postman
	env     staff.Env
}

// NewStaffRouter creates a new instance of StaffController
func NewStaffRouter(db *sqlx.DB, p postoffice.Postman) StaffRouter {
	return StaffRouter{
		env:     staff.Env{DB: db},
		postman: p,
	}
}

// Login verifies a user's user name and password. Headers: `X-User-Ip`.
func (router StaffRouter) Login(w http.ResponseWriter, req *http.Request) {
	var login employee.Login

	// `400 Bad Request` if body content cannot be parsed as JSON
	if err := gorest.ParseJSON(req.Body, &login); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	login.Sanitize()

	account, err := router.env.Login(login)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	userIP := req.Header.Get("X-User-Ip")
	go func() {
		err := router.env.UpdateLastLogin(login, userIP)
		if err != nil {
			logger.WithField("trace", "StaffRouter.Login")
		}
	}()

	go func() {
		parcel, err := account.SignUpParcel()
		if err != nil {
			logger.WithField("trace", "StaffRouter.Login").Error(err)
		}
		if err := router.postman.Deliver(parcel); err != nil {
			logger.WithField("trace", "StaffRouter.Login").Error(err)
		}
	}()

	// `200 OK`
	view.Render(w, view.NewResponse().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST /staff/password-reset/letter
//
// Input {email: string}
func (router StaffRouter) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	var th employee.TokenHolder
	if err := gorest.ParseJSON(req.Body, &th); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	th.Sanitize()
	if r := th.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}
	if err := th.GenerateToken(); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	profile, err := router.env.Load(staff.ColumnEmail, th.Email)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}
	if !profile.IsActive {
		view.Render(w, view.NewNotFound())
		return
	}

	if err := router.env.SavePwResetToken(th); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	parcel, err := profile.PasswordResetParcel(th.Token)
	if err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	go func() {
		if err := router.postman.Deliver(parcel); err != nil {
			logger.WithField("trace", "StaffRouter.ForgotPassword")
		}
	}()

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
//
// 	GET /password-reset/tokens/{token}
//
// Output {email: string}
func (router StaffRouter) VerifyToken(w http.ResponseWriter, req *http.Request) {
	token, err := GetURLParam(req, "token").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	th, err := router.env.LoadResetToken(token)

	// `404 Not Found`
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `200 OK`
	view.Render(w, view.NewResponse().SetBody(th))
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
//
//	POST /password-reset
//
// Input {token: string, password: string}
func (router StaffRouter) ResetPassword(w http.ResponseWriter, req *http.Request) {
	var reset employee.PasswordReset

	// `400 Bad Request`
	if err := gorest.ParseJSON(req.Body, &reset); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	reset.Sanitize()
	if r := reset.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	if err := router.env.ResetPassword(reset); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}

// Creates creates a new account.
func (router StaffRouter) Create(w http.ResponseWriter, req *http.Request) {
	var a employee.Account

	if err := gorest.ParseJSON(req.Body, &a); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := a.GenerateID(); err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}
	if err := a.GeneratePassword(); err != nil {
		view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	if r := a.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	if err := router.env.Create(a); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(a))
}

// List shows all staff.
func (router StaffRouter) List(w http.ResponseWriter, req *http.Request) {

	pagination := gorest.GetPagination(req)

	profiles, err := router.env.List(pagination)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(profiles))
}

// Profile shows a staff's profile.
//
//	 GET /staff/{id}
func (router StaffRouter) Profile(w http.ResponseWriter, req *http.Request) {

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p, err := router.env.Load(staff.ColumnStaffId, id)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(p))
}

func (router StaffRouter) Update(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var p employee.Profile
	if err := gorest.ParseJSON(req.Body, &p); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if r := p.Validate(); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	p.ID = id
	if err := router.env.Update(p); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

func (router StaffRouter) Delete(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var vip = struct {
		Revoke bool `json:"revokeVip"`
	}{}
	if err := gorest.ParseJSON(req.Body, &vip); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := router.env.Deactivate(id, vip.Revoke); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

func (router StaffRouter) Reinstate(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := router.env.Activate(id); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewNoContent())
}

// UpdateDisplayName lets user to change displayed name.
//
//	PATCH /staff/display-name
//
// Input {displayName: string}
func (router StaffRouter) UpdateDisplayName(w http.ResponseWriter, req *http.Request) {

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	displayName, err := GetString(req.Body, "displayName")
	// `400 Bad Request`
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	// `422 Unprocessable Entity`
	if r := util.RequireNotEmptyWithMax(displayName, 255, "displayName"); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	err = router.env.UpdateDisplayName(displayName, id)

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

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	email, err := GetString(req.Body, "email")
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	// `422 Unprocessable Entity`
	if r := util.RequireEmail(email); r != nil {
		view.Render(w, view.NewUnprocessable(r))
		return
	}

	err = router.env.UpdateEmail(email, id)

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

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var p employee.Password
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

	// `403 Forbidden` if old password is wrong
	if err := router.env.UpdatePassword(p, id); err != nil {
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
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var login reader.Login
	// `400 Bad Request` for invalid JSON
	if err := gorest.ParseJSON(req.Body, &login); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	login.Sanitize()
	ftcAccount, err := router.env.MyftAuth(login)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	err = router.env.LinkFtc(employee.Myft{
		StaffID: id,
		MyftID:  ftcAccount.ID,
	})

	// `404 Not Found` if `email` + `password` verification failed.
	if err != nil {
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
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	ftcAccounts, err := router.env.ListMyft(id)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `200 OK`
	view.Render(w, view.NewResponse().NoCache().SetBody(ftcAccounts))
}

// DeleteMyft deletes an ftc account owned by current user.
//
//	DELETE /staff/myft
//
// Input {ftcId: string}
func (router StaffRouter) DeleteMyft(w http.ResponseWriter, req *http.Request) {

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var myft employee.Myft
	if err := gorest.ParseJSON(req.Body, &myft); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}
	myft.StaffID = id

	if err := router.env.UnlinkFtc(myft); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	view.Render(w, view.NewNoContent())
}
