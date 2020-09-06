package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateAccount handles creating a sandbox account.
// Input: reader.SandboxInput
// email: string
// password: string
func (router ReaderRouter) CreateSandboxUser(c echo.Context) error {
	claims := getPassportClaims(c)

	var input reader.SandboxInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account := reader.NewSandboxFtcAccount(input, claims.Username)

	err := router.readerRepo.CreateSandboxUser(account)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// ListUsers retrieves all sandbox user.
func (router ReaderRouter) ListSandboxUsers(c echo.Context) error {
	users, err := router.readerRepo.ListSandboxFtcAccount()
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, users)
}

// LoadAccount loads a sandbox reader with membership.
func (router ReaderRouter) LoadSandboxAccount(c echo.Context) error {
	id := c.Param("id")
	account, err := router.readerRepo.LoadSandboxAccount(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// ChangeSandboxPassword overrides current password.
// Input:
// ftcId: string;
// password: string;
func (router ReaderRouter) ChangeSandboxPassword(c echo.Context) error {
	var input reader.SandboxFtcAccount
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	err := router.readerRepo.ChangePassword(input)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
