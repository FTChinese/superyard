package controller

import (
	"database/sql"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/staff"
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
func (router StaffRouter) Login(c echo.Context) error {
	var login employee.Login

	// `400 Bad Request` if body content cannot be parsed as JSON
	if err := c.Bind(&login); err != nil {
		return util.NewBadRequest(err.Error())
	}

	login.Sanitize()

	if ie := login.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}

	account, err := router.env.Login(login)
	if err != nil {
		return util.NewDBFailure(err)
	}

	userIP := c.RealIP()
	go func() {
		err := router.env.UpdateLastLogin(login, userIP)
		if err != nil {
			logger.WithField("trace", "StaffRouter.Login").Error(err)
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
	return c.JSON(http.StatusOK, account)
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST /staff/password-reset/letter
//
// Input {email: string}
func (router StaffRouter) ForgotPassword(c echo.Context) error {

	var th employee.TokenHolder
	if err := c.Bind(&th); err != nil {
		return util.NewBadRequest(err.Error())
	}

	th.Sanitize()
	if ie := th.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}
	if err := th.GenerateToken(); err != nil {
		return err
	}

	account, err := router.env.RetrieveAccount(employee.ColumnEmail, th.Email)
	if err != nil {
		return util.NewDBFailure(err)
	}
	if !account.IsActive {
		return util.NewNotFound("")
	}

	if err := router.env.SavePwResetToken(th); err != nil {
		return util.NewDBFailure(err)
	}

	parcel, err := account.PasswordResetParcel(th.Token)
	if err != nil {
		return err
	}

	go func() {
		if err := router.postman.Deliver(parcel); err != nil {
			logger.WithField("trace", "StaffRouter.ForgotPassword")
		}
	}()

	// `204 No Content`
	return c.NoContent(http.StatusNoContent)
}

// VerifyToken checks if a token exists when user clicked the link in password reset letter
//
// 	GET /password-reset/tokens/{token}
//
// Output {email: string}
func (router StaffRouter) VerifyToken(c echo.Context) error {
	token := c.Param("token")

	th, err := router.env.LoadResetToken(token)

	// `404 Not Found`
	if err != nil {
		return util.NewDBFailure(err)
	}

	// `200 OK`
	return c.JSON(http.StatusOK, th)
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
//
//	POST /password-reset
//
// Input {token: string, password: string}
func (router StaffRouter) ResetPassword(c echo.Context) error {
	var reset employee.PasswordReset

	// `400 Bad Request`
	if err := c.Bind(&reset); err != nil {
		return util.NewBadRequest(err.Error())
	}

	reset.Sanitize()
	if ie := reset.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}

	th, err := router.env.LoadResetToken(reset.Token)
	if err != nil {
		return util.NewDBFailure(err)
	}

	account, err := router.env.RetrieveAccount(employee.ColumnEmail, th.Email)
	if err != nil {
		return util.NewDBFailure(err)
	}
	account.Password = reset.Password

	if err := router.env.UpdatePassword(account); err != nil {
		return util.NewDBFailure(err)
	}

	if err := router.env.DeleteResetToken(reset.Token); err != nil {
		return util.NewDBFailure(err)
	}

	// `204 No Content`
	return c.NoContent(http.StatusNoContent)
}

// Creates creates a new account.
func (router StaffRouter) Create(c echo.Context) error {
	var a employee.Account

	if err := c.Bind(&a); err != nil {
		return util.NewBadRequest(err.Error())
	}

	a.GenerateID()

	if err := a.GeneratePassword(); err != nil {
		return util.NewBadRequest(err.Error())
	}

	if ie := a.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}

	if err := router.env.Create(a); err != nil {
		return util.NewDBFailure(err)
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

	return c.JSON(http.StatusOK, a)
}

// ListStaff shows all staff.
func (router StaffRouter) List(c echo.Context) error {

	var pagination util.Pagination
	if err := c.Bind(&pagination); err != nil {
		return util.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	logger.Infof("Pagination: %+v", pagination)

	profiles, err := router.env.ListStaff(pagination)
	if err != nil {
		return util.NewDBFailure(err)
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
			_ = router.env.AddID(p.Account)
		}
	}()

	return c.JSON(http.StatusOK, profiles)
}

// Profile shows a staff's profile.
//
//	 GET /staff/{id}
func (router StaffRouter) Profile(c echo.Context) error {

	id := c.Param("id")

	logger.Infof("Profile for staff: %s", id)

	p, err := router.env.RetrieveProfile(id)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, p)
}

func (router StaffRouter) Update(c echo.Context) error {
	log := logger.WithField("trace", "StaffRouter.Update")

	id := c.Param("id")

	p, err := router.env.RetrieveProfile(id)
	if err != nil {
		return util.NewDBFailure(err)
	}

	if err := c.Bind(&p); err != nil {
		log.Error(err)
		return util.NewBadRequest(err.Error())
	}

	if ie := p.Validate(); ie != nil {
		log.Error(ie)
		return util.NewUnprocessable(ie)
	}
	// In case input data contains id field.
	p.ID = null.StringFrom(id)

	if err := router.env.UpdateProfile(p); err != nil {
		log.Error(err)
		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router StaffRouter) Delete(c echo.Context) error {
	id := c.Param("id")

	var vip = struct {
		Revoke bool `json:"revokeVip"`
	}{}
	if err := c.Bind(&vip); err != nil {
		return util.NewBadRequest(err.Error())
	}

	if err := router.env.Deactivate(id); err != nil {
		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router StaffRouter) Reinstate(c echo.Context) error {
	id := c.Param("id")

	if err := router.env.Activate(id); err != nil {
		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdatePassword lets user to change password.
//
//	PATCH /staff/{id}/password
//
// Input {oldPassword: string, newPassword: string}
func (router StaffRouter) UpdatePassword(c echo.Context) error {

	id := c.Param("id")

	var p employee.Password
	if err := c.Bind(&p); err != nil {
		return util.NewBadRequest(err.Error())
	}
	p.Sanitize()

	// `422 Unprocessable Entity`
	if ie := p.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}

	account, err := router.env.VerifyPassword(employee.Account{
		ID:       null.StringFrom(id),
		Password: p.Old,
	})
	if err != nil {
		// No rows means password is incorrect.
		if err == sql.ErrNoRows {
			return util.NewForbidden("Current password incorrect")
		}
		return util.NewDBFailure(err)
	}
	account.Password = p.New

	// `403 Forbidden` if old password is wrong
	if err := router.env.UpdatePassword(account); err != nil {
		return util.NewDBFailure(err)
	}

	// `204 No Content`
	return c.NoContent(http.StatusNoContent)
}

// Search finds an employee.
// Query parameter: q=<user name>
func (router StaffRouter) Search(c echo.Context) error {
	q := c.QueryParam("q")
	if q == "" {
		return util.NewBadRequest("Missing query parameter q")
	}

	account, err := router.env.RetrieveAccount(employee.ColumnUserName, q)
	if err != nil {
		return util.NewDBFailure(err)
	}

	if account.ID.IsZero() {
		account.GenerateID()
		go func() {
			_ = router.env.AddID(account)
		}()
	}

	return c.JSON(http.StatusOK, account)
}
