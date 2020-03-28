package controller

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/user"
	"net/http"
)

type UserRouter struct {
	repo    user.Env
	postman postoffice.Postman
}

func NewUserRouter(db *sqlx.DB, p postoffice.Postman) UserRouter {
	return UserRouter{
		repo:    user.Env{DB: db},
		postman: p,
	}
}

func (router UserRouter) Login(c echo.Context) error {
	var login employee.Login

	// `400 Bad Request` if body content cannot be parsed as JSON
	if err := c.Bind(&login); err != nil {
		return render.NewBadRequest(err.Error())
	}

	login.Sanitize()
	if ve := login.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.Login(login)
	if err != nil {
		return render.NewDBError(err)
	}

	userIP := c.RealIP()
	go func() {
		_ = router.repo.UpdateLastLogin(login, userIP)
	}()

	if account.ID.IsZero() {
		account.ID = null.StringFrom(employee.GenStaffID())
		go func() {
			_ = router.repo.AddID(account)
		}()
	}

	// Includes JWT in response.
	jwtAccount, err := employee.NewJWTAccount(account)
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
// Input {email: string}
func (router UserRouter) ForgotPassword(c echo.Context) error {
	var pr employee.PasswordReset
	if err := c.Bind(&pr); err != nil {
		return render.NewBadRequest(err.Error())
	}

	pr.Sanitize()
	if ve := pr.ValidateEmail(); ve != nil {
		return render.NewUnprocessable(ve)
	}
	if err := pr.GenerateToken(); err != nil {
		return err
	}

	account, err := router.repo.AccountByEmail(pr.Email)
	if err != nil {
		return render.NewDBError(err)
	}
	if !account.IsActive {
		return render.NewNotFound("")
	}
	// Add id is missing.
	if account.ID.IsZero() {
		account.ID = null.StringFrom(employee.GenStaffID())
		go func() {
			_ = router.repo.AddID(account)
		}()
	}

	// Generate token
	err = pr.GenerateToken()
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	// Save toke and email
	if err := router.repo.SavePwResetToken(pr); err != nil {
		return render.NewDBError(err)
	}

	// Create email content
	parcel, err := account.PasswordResetParcel(pr.Token)
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

// VerifyToken checks if a token exists when user clicked the link in password reset letter
//
// 	GET /password-reset/tokens/{token}
//
// Output employee.Account.
func (router UserRouter) VerifyToken(c echo.Context) error {
	token := c.Param("token")

	account, err := router.repo.AccountByResetToken(token)

	// `404 Not Found`
	if err != nil {
		return render.NewDBError(err)
	}

	// `200 OK`
	return c.JSON(http.StatusOK, account)
}

// ResetPassword verifies password reset token and allows user to submit new password if the token is valid
//
//	POST /password-reset
//
// Input {token: string, password: string}
func (router UserRouter) ResetPassword(c echo.Context) error {
	var reset employee.PasswordReset

	// `400 Bad Request`
	if err := c.Bind(&reset); err != nil {
		return render.NewBadRequest(err.Error())
	}

	reset.Sanitize()
	if ve := reset.ValidatePass(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByResetToken(reset.Token)
	if err != nil {
		return render.NewDBError(err)
	}

	// Change password.
	err = router.repo.UpdatePassword(account.Credentials(reset.Password))

	if err != nil {
		return render.NewDBError(err)
	}

	if err := router.repo.DisableResetToken(reset.Token); err != nil {
		return render.NewDBError(err)
	}

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

	var ba employee.BaseAccount
	if err := c.Bind(ba); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := ba.ValidateEmail(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByID(claims.StaffID)
	if err != nil {
		return render.NewDBError(err)
	}

	// If email is not changed, do not touch db.
	if account.Email == ba.Email {
		return c.NoContent(http.StatusNoContent)
	}

	// Update the account instance.
	if ve := account.SetEmail(ba.Email); ve != nil {
		return render.NewUnprocessable(ve)
	}

	err = router.repo.SetEmail(account)
	if err != nil {
		if util.IsAlreadyExists(err) {
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

	var ba employee.BaseAccount
	if err := c.Bind(ba); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := ba.ValidateDisplayName(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByID(claims.StaffID)
	if err != nil {
		return render.NewDBError(err)
	}

	// If display name is not changed, do not touch db.
	if account.DisplayName == ba.DisplayName {
		return c.NoContent(http.StatusNoContent)
	}

	// Update the account instance.
	account.DisplayName = ba.DisplayName

	err = router.repo.UpdateDisplayName(account)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdatePassword lets user to change password.
//
// Input {oldPassword: string, newPassword: string}
func (router UserRouter) ChangePassword(c echo.Context) error {
	claims := getAccountClaims(c)

	var p employee.Password
	if err := c.Bind(&p); err != nil {
		return render.NewBadRequest(err.Error())
	}
	p.Sanitize()

	// `422 Unprocessable Entity`
	if ve := p.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// Verify old password.
	account, err := router.repo.VerifyPassword(employee.Credentials{
		ID: claims.StaffID,
		Login: employee.Login{
			Password: p.Old,
		},
	})

	if err != nil {
		return render.NewDBError(err)
	}

	err = router.repo.UpdatePassword(account.Credentials(p.New))
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
