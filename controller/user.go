package controller

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/pkg/db"
	"gitlab.com/ftchinese/superyard/pkg/staff"
	"gitlab.com/ftchinese/superyard/repository/user"
	"net/http"
)

type UserRouter struct {
	repo    user.Env
	postman postoffice.PostOffice
}

func NewUserRouter(db *sqlx.DB, p postoffice.PostOffice) UserRouter {
	return UserRouter{
		repo:    user.Env{DB: db},
		postman: p,
	}
}

// Login verifies user name and password.
//
// Input:
// {userName: string, password: string}
func (router UserRouter) Login(c echo.Context) error {
	var input staff.InputData

	// `400 Bad Request` if body content cannot be parsed as JSON
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateLogin(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	login := input.Login()
	account, err := router.repo.Login(login)
	if err != nil {
		return render.NewDBError(err)
	}

	userIP := c.RealIP()
	go func() {
		_ = router.repo.UpdateLastLogin(login, userIP)
	}()

	if account.ID.IsZero() {
		account.ID = null.StringFrom(staff.GenStaffID())
		go func() {
			_ = router.repo.AddID(account)
		}()
	}

	// Includes JWT in response.
	jwtAccount, err := staff.NewJWTAccount(account)
	if err != nil {
		return render.NewUnauthorized(err.Error())
	}

	// `200 OK`
	return c.JSON(http.StatusOK, jwtAccount)
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST /adminRepo/password-reset/letter
//
// Input {email: string, sourceUrl?: string}
func (router UserRouter) ForgotPassword(c echo.Context) error {
	var input staff.InputData
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateEmail(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	session, err := staff.NewPwResetSession(input.Email)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	account, err := router.repo.AccountByEmail(session.Email)
	if err != nil {
		return render.NewDBError(err)
	}
	if !account.IsActive {
		return render.NewNotFound("")
	}
	// Add id is missing.
	if account.ID.IsZero() {
		account.ID = null.StringFrom(staff.GenStaffID())
		go func() {
			_ = router.repo.AddID(account)
		}()
	}

	// Save toke and email
	if err := router.repo.SavePwResetSession(session); err != nil {
		return render.NewDBError(err)
	}

	// CreateStaff email content
	parcel, err := account.PasswordResetParcel(session)
	if err != nil {
		return err
	}

	// Send email
	go func() {
		if err := router.postman.Deliver(parcel); err != nil {
			logger.WithField("trace", "UserRouter.ForgotPassword")
		}
	}()

	// `204 No Content`
	return c.NoContent(http.StatusNoContent)
}

// VerifyResetToken checks if a token exists when user clicked the link in password reset letter
//
// 	GET /password-reset/tokens/{token}
//
// Output employee.Account.
func (router UserRouter) VerifyResetToken(c echo.Context) error {
	token := c.Param("token")

	session, err := router.repo.LoadPwResetSession(token)
	if err != nil {
		return render.NewDBError(err)
	}

	if session.IsUsed || session.IsExpired() {
		return render.NewNotFound("token already used or expired")
	}

	// `200 OK`
	return c.JSON(http.StatusOK, map[string]string{
		"email": session.Email,
	})
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
//
//	POST /password-reset
//
// Input {token: string, password: string}
func (router UserRouter) ResetPassword(c echo.Context) error {
	var input staff.InputData

	// `400 Bad Request`
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidatePasswordReset(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByResetToken(input.Token)
	if err != nil {
		return render.NewDBError(err)
	}

	// Change password.
	err = router.repo.UpdatePassword(staff.Credentials{
		UserName: account.UserName,
		Password: input.Password,
	})

	if err != nil {
		return render.NewDBError(err)
	}

	go func() {
		_ = router.repo.DisableResetToken(input.Token)
	}()

	// `204 No Content`
	return c.NoContent(http.StatusNoContent)
}

// Account returns a logged in user's account.
func (router UserRouter) Account(c echo.Context) error {
	claims := getAccountClaims(c)

	account, err := router.repo.AccountByID(claims.StaffID)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// SetEmail set the email column.
// Input: {email: string}
func (router UserRouter) SetEmail(c echo.Context) error {
	claims := getAccountClaims(c)

	var input staff.InputData
	if err := c.Bind(input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateEmail(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByID(claims.StaffID)
	if err != nil {
		return render.NewDBError(err)
	}

	// If email is not changed, do not touch db.
	if account.Email == input.Email {
		return c.NoContent(http.StatusNoContent)
	}

	account.Email = input.Email

	err = router.repo.SetEmail(account)
	if err != nil {
		if db.IsAlreadyExists(err) {
			return render.NewAlreadyExists("email")
		}
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeDisplayName updates display name.
// Input {displayName: string}
func (router UserRouter) ChangeDisplayName(c echo.Context) error {
	claims := getAccountClaims(c)

	var input staff.InputData
	if err := c.Bind(input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateDisplayName(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByID(claims.StaffID)
	if err != nil {
		return render.NewDBError(err)
	}

	// If display name is not changed, do not touch db.
	if account.DisplayName == input.DisplayName {
		return c.NoContent(http.StatusNoContent)
	}

	// Update the account instance.
	account.DisplayName = input.DisplayName

	err = router.repo.UpdateDisplayName(account)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdatePassword lets user to change password.
//
// Input {oldPassword: string, password: string}
func (router UserRouter) UpdatePassword(c echo.Context) error {
	claims := getAccountClaims(c)

	var input staff.InputData
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	// `422 Unprocessable Entity`
	if ve := input.ValidatePwUpdater(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// Verify old password.
	account, err := router.repo.VerifyPassword(staff.PasswordVerifier{
		StaffID:     claims.StaffID,
		OldPassword: input.OldPassword,
	})

	if err != nil {
		return render.NewDBError(err)
	}

	err = router.repo.UpdatePassword(staff.Credentials{
		UserName: account.UserName,
		Password: input.Password,
	})
	if err != nil {
		return render.NewDBError(err)
	}

	// `204 No Content`
	return c.NoContent(http.StatusNoContent)
}

// Profile shows a adminRepo's profile.
func (router UserRouter) Profile(c echo.Context) error {
	claims := getAccountClaims(c)

	p, err := router.repo.RetrieveProfile(claims.StaffID)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, p)
}
