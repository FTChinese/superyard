package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/employee"
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
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	login.Sanitize()

	account, err := router.env.Login(login)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	userIP := req.Header.Get("X-User-Ip")
	go func() {
		err := router.env.UpdateLastLogin(login, userIP)
		if err != nil {
			logger.WithField("trace", "StaffRouter.Login")
		}
	}()

	if account.ID.IsZero() {
		account.GenerateID()
		go func() {
			if err := router.env.AddID(account); err != nil {
				logger.WithField("trace", "Env.Login").Error(err)
			}
		}()
	}
	// `200 OK`
	_ = view.Render(w, view.NewResponse().SetBody(account))
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST /staff/password-reset/letter
//
// Input {email: string}
func (router StaffRouter) ForgotPassword(w http.ResponseWriter, req *http.Request) {

	var th employee.TokenHolder
	if err := gorest.ParseJSON(req.Body, &th); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	th.Sanitize()
	if r := th.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}
	if err := th.GenerateToken(); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	profile, err := router.env.Load(staff.ColumnEmail, th.Email)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}
	if !profile.IsActive {
		_ = view.Render(w, view.NewNotFound())
		return
	}

	if err := router.env.SavePwResetToken(th); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	parcel, err := profile.PasswordResetParcel(th.Token)
	if err != nil {
		_ = view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	go func() {
		if err := router.postman.Deliver(parcel); err != nil {
			logger.WithField("trace", "StaffRouter.ForgotPassword")
		}
	}()

	// `204 No Content`
	_ = view.Render(w, view.NewNoContent())
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
//
// 	GET /password-reset/tokens/{token}
//
// Output {email: string}
func (router StaffRouter) VerifyToken(w http.ResponseWriter, req *http.Request) {
	token, err := GetURLParam(req, "token").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	th, err := router.env.LoadResetToken(token)

	// `404 Not Found`
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	// `200 OK`
	_ = view.Render(w, view.NewResponse().SetBody(th))
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
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	reset.Sanitize()
	if r := reset.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	if err := router.env.ResetPassword(reset); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	_ = view.Render(w, view.NewNoContent())
}

// Creates creates a new account.
func (router StaffRouter) Create(w http.ResponseWriter, req *http.Request) {
	var a employee.Account

	if err := gorest.ParseJSON(req.Body, &a); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	a.GenerateID()

	if err := a.GeneratePassword(); err != nil {
		_ = view.Render(w, view.NewInternalError(err.Error()))
		return
	}

	if r := a.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	if err := router.env.Create(a); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	go func() {
		parcel, err := a.SignUpParcel()
		if err != nil {
			logger.WithField("trace", "StaffRouter.Login").Error(err)
		}
		if err := router.postman.Deliver(parcel); err != nil {
			logger.WithField("trace", "StaffRouter.Login").Error(err)
		}
	}()

	_ = view.Render(w, view.NewResponse().SetBody(a))
}

// List shows all staff.
func (router StaffRouter) List(w http.ResponseWriter, req *http.Request) {

	pagination := gorest.GetPagination(req)

	profiles, err := router.env.List(pagination)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(profiles))
}

// Profile shows a staff's profile.
//
//	 GET /staff/{id}
func (router StaffRouter) Profile(w http.ResponseWriter, req *http.Request) {

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p, err := router.env.Load(staff.ColumnStaffId, id)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(p))
}

func (router StaffRouter) Update(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var p employee.Profile
	if err := gorest.ParseJSON(req.Body, &p); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if r := p.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	p.ID = null.StringFrom(id)

	if err := router.env.Update(p); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

func (router StaffRouter) Delete(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var vip = struct {
		Revoke bool `json:"revokeVip"`
	}{}
	if err := gorest.ParseJSON(req.Body, &vip); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := router.env.Deactivate(id, vip.Revoke); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

func (router StaffRouter) Reinstate(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := router.env.Activate(id); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}

// UpdateDisplayName lets user to change displayed name.
//
//	PATCH /staff/display-name
//
// Input {displayName: string}
func (router StaffRouter) UpdateDisplayName(w http.ResponseWriter, req *http.Request) {

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	displayName, err := GetString(req.Body, "displayName")
	// `400 Bad Request`
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	// `422 Unprocessable Entity`
	if r := util.RequireNotEmptyWithMax(displayName, 255, "displayName"); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	err = router.env.UpdateDisplayName(displayName, id)

	// `422 Unprocessable Entity` if this `displayName` already exists
	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Code = view.CodeAlreadyExists
			reason.Field = "displayName"
			_ = view.Render(w, view.NewUnprocessable(reason))
			return
		}
		_ = view.Render(w, view.NewDBFailure(err))

		return
	}

	// `204 No Content`
	_ = view.Render(w, view.NewNoContent())
}

// UpdateEmail lets user to change email.
//
//	PATCH /staff/email
//
// Input {email: string}
func (router StaffRouter) UpdateEmail(w http.ResponseWriter, req *http.Request) {

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	email, err := GetString(req.Body, "email")
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	// `422 Unprocessable Entity`
	if r := util.RequireEmail(email); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	err = router.env.UpdateEmail(email, id)

	// `422 Unprocessable Entity`
	if err != nil {
		if IsAlreadyExists(err) {
			reason := view.NewReason()
			reason.Field = "email"
			reason.Code = view.CodeAlreadyExists
			_ = view.Render(w, view.NewUnprocessable(reason))
			return
		}
		_ = view.Render(w, view.NewDBFailure(err))

		return
	}

	// `204 No Content` if updated successfully.
	_ = view.Render(w, view.NewNoContent())
}

// UpdatePassword lets user to change password.
//
//	PATCH /staff/password
//
// Input {oldPassword: string, newPassword: string}
func (router StaffRouter) UpdatePassword(w http.ResponseWriter, req *http.Request) {

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var p employee.Password
	if err := gorest.ParseJSON(req.Body, &p); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}
	p.Sanitize()

	// `422 Unprocessable Entity`
	if r := p.Validate(); r != nil {
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	// `403 Forbidden` if old password is wrong
	if err := router.env.UpdatePassword(p, id); err != nil {
		if err == util.ErrWrongPassword {
			_ = view.Render(w, view.NewForbidden(err.Error()))
			return
		}
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	_ = view.Render(w, view.NewNoContent())
}
