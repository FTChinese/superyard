package controller

import (
	"database/sql"
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
		_ = view.JSON(w, view.NewBadRequest(err.Error()))
		return
	}

	login.Sanitize()

	if r := login.Validate(); r != nil {
		_ = view.JSON(w, view.NewUnprocessable(r))
		return
	}

	account, err := router.env.Login(login)
	if err != nil {
		_ = view.JSON(w, view.NewDBFailure(err))
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
	_ = view.JSON(w, view.NewResponse().SetBody(account))
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

	account, err := router.env.RetrieveAccount(employee.ColumnEmail, th.Email)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}
	if !account.IsActive {
		_ = view.Render(w, view.NewNotFound())
		return
	}

	if err := router.env.SavePwResetToken(th); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	parcel, err := account.PasswordResetParcel(th.Token)
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

	th, err := router.env.LoadResetToken(reset.Token)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	account, err := router.env.RetrieveAccount(employee.ColumnEmail, th.Email)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}
	account.Password = reset.Password

	if err := router.env.UpdatePassword(account); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	if err := router.env.DeleteResetToken(reset.Token); err != nil {
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

// ListStaff shows all staff.
func (router StaffRouter) List(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	logger.Infof("Pagination: %+v", pagination)

	profiles, err := router.env.ListStaff(pagination)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	var noIDs []employee.Profile

	for i, p := range profiles {
		if p.ID.IsZero() {
			// NOTE: the element in the original slice is not
			// touched! Replace the older one with the new one
			p.GenerateID()
			profiles[i] = p
			noIDs = append(noIDs, p)
		}
	}

	go func() {
		for _, p := range noIDs {
			if err := router.env.AddID(p.Account); err != nil {
				logger.WithField("trace", "StaffRouter.ListStaff")
			}
		}
	}()

	_ = view.Render(w, view.NewResponse().SetBody(profiles))
}

// Profile shows a staff's profile.
//
//	 GET /staff/{id}
func (router StaffRouter) Profile(w http.ResponseWriter, req *http.Request) {

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		logger.WithField("trace", "StaffRouter.Profile").Error(err)
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	logger.Infof("Profile for staff: %s", id)

	p, err := router.env.RetrieveProfile(id)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(p))
}

func (router StaffRouter) Update(w http.ResponseWriter, req *http.Request) {
	log := logger.WithField("trace", "StaffRouter.Update")

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		log.Error(err)
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p, err := router.env.RetrieveProfile(id)
	if err != nil {
		log.Error(err)
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	if err := gorest.ParseJSON(req.Body, &p); err != nil {
		log.Error(err)
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if r := p.Validate(); r != nil {
		log.Error(r)
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}
	// In case input data contains id field.
	p.ID = null.StringFrom(id)

	if err := router.env.UpdateProfile(p); err != nil {
		log.Error(err)
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

	if err := router.env.Deactivate(id); err != nil {
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

// UpdatePassword lets user to change password.
//
//	PATCH /staff/{id}/password
//
// Input {oldPassword: string, newPassword: string}
func (router StaffRouter) UpdatePassword(w http.ResponseWriter, req *http.Request) {

	log := logger.WithField("trace", "StaffRouter.UpdatePassword")

	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		log.Error(err)
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var p employee.Password
	if err := gorest.ParseJSON(req.Body, &p); err != nil {
		log.Error(err)
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}
	p.Sanitize()

	// `422 Unprocessable Entity`
	if r := p.Validate(); r != nil {
		log.Error(err)
		_ = view.Render(w, view.NewUnprocessable(r))
		return
	}

	account, err := router.env.VerifyPassword(employee.Account{
		ID:       null.StringFrom(id),
		Password: p.Old,
	})
	if err != nil {
		log.Error(err)
		// No rows means password is incorrect.
		if err == sql.ErrNoRows {
			_ = view.Render(w, view.NewForbidden(util.ErrWrongPassword.Error()))
			return
		}
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}
	account.Password = p.New

	// `403 Forbidden` if old password is wrong
	if err := router.env.UpdatePassword(account); err != nil {
		log.Error(err)
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	// `204 No Content`
	_ = view.Render(w, view.NewNoContent())
}
