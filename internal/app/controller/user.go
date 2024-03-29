package controller

import (
	"net/http"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/auth"
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/letter"
	"github.com/FTChinese/superyard/pkg/postman"
	"github.com/labstack/echo/v4"
)

type UserRouter struct {
	guard   AuthGuard
	repo    auth.Env
	postman postman.Postman
}

func NewUserRouter(gormDBs db.MultiGormDBs, p postman.Postman, g AuthGuard) UserRouter {
	return UserRouter{
		guard:   g,
		repo:    auth.NewEnv(gormDBs),
		postman: p,
	}
}

// Login verifies user name and password.
//
// Input: {userName: string, password: string}
// Response
// 400 if request body cannot be parsed;
// 422 if validation failed. The response JSON has a `error` field:
// { field: "userName", code: "missing_field"} if userName is missing;
// { field: "userName", code: "invalid"} if userName exceeds 64 chars;
// { field: "password", code: "missing_field"} if password is missing;
// { field: "password", code: "invalid"} if password exceeds 64 chars.
func (router UserRouter) Login(c echo.Context) error {
	var input user.Credentials

	// `400 Bad Request` if body content cannot be parsed as JSON
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.Login(input)
	if err != nil {
		return render.NewDBError(err)
	}

	// Includes JWT in response.
	passport, err := user.NewPassport(account, router.guard.signingKey)

	if err != nil {
		return render.NewUnauthorized(err.Error())
	}

	// `200 OK`
	return c.JSON(http.StatusOK, passport)
}

// ForgotPassword checks user's email and send a password reset letter if it is valid
//
//	POST /password-reset/letter
//
// Input {email: string, sourceUrl?: string}
// Response:
// 204 if everything is ok.
func (router UserRouter) ForgotPassword(c echo.Context) error {
	var input user.ParamsForgotPassLetter
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	session, err := user.NewPwResetSession(input.Email)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	account, err := router.repo.AccountByEmail(session.Email)
	if err != nil {
		return render.NewDBError(err)
	}

	// Save toke and email
	if err := router.repo.SavePwResetSession(session); err != nil {
		return render.NewDBError(err)
	}

	// CreateStaff email content
	parcel, err := letter.PasswordResetParcel(account, session.BuildURL(input.SourceURL))
	if err != nil {
		return err
	}

	// Send email
	go func() {
		if err := router.postman.Deliver(parcel); err != nil {

		}
	}()

	// `204 No Content`
	return c.NoContent(http.StatusNoContent)
}

// VerifyResetToken checks if a token exists when user clicked the link in password reset letter
//
//	GET /password-reset/tokens/{token}
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
	var input user.ParamsResetPass

	// `400 Bad Request`
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	session, err := router.repo.LoadPwResetSession(input.Token)
	if err != nil {
		return render.NewDBError(err)
	}

	account, err := router.repo.AccountByEmail(session.Email)
	if err != nil {
		return render.NewDBError(err)
	}

	// Change password.
	err = router.repo.UpdatePassword(user.Credentials{
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
	claims := getPassportClaims(c)

	account, err := router.repo.AccountByID(claims.UserID)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// SetEmail set the email column.
// Input: {email: string}
func (router UserRouter) SetEmail(c echo.Context) error {
	claims := getPassportClaims(c)

	var input user.ParamsEmail
	if err := c.Bind(input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByID(claims.UserID)
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
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeDisplayName updates display name.
// Input {displayName: string}
func (router UserRouter) ChangeDisplayName(c echo.Context) error {
	claims := getPassportClaims(c)

	var input user.ParamsDisplayName
	if err := c.Bind(input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account, err := router.repo.AccountByID(claims.UserID)
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
	claims := getPassportClaims(c)

	var input user.ParamsPasswords
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	// `422 Unprocessable Entity`
	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// Verify old password.
	account, err := router.repo.VerifyPassword(claims.UserID, input.OldPassword)

	if err != nil {
		return render.NewDBError(err)
	}

	err = router.repo.UpdatePassword(user.Credentials{
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
	claims := getPassportClaims(c)

	p, err := router.repo.RetrieveProfile(claims.UserID)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, p)
}
