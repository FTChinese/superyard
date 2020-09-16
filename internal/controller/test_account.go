package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateAccount handles creating a sandbox account.
//
// POST /sandbox
//
// Input: reader.TestAccountInput
// email: string
// password: string
func (router ReaderRouter) CreateTestUser(c echo.Context) error {
	claims := getPassportClaims(c)

	var input reader.TestAccountInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	account := reader.NewTestFtcAccount(input, claims.Username)

	err := router.readerRepo.CreateTestUser(account)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// ListUsers retrieves all sandbox user.
//
// GET /sandbox
func (router ReaderRouter) ListTestUsers(c echo.Context) error {
	var p gorest.Pagination
	if err := c.Bind(&p); err != nil {
		return render.NewBadRequest(err.Error())
	}

	users, err := router.readerRepo.ListTestFtcAccount(p)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, users)
}

// LoadAccount loads a sandbox reader with membership.
//
// GET /sandbox/:id
func (router ReaderRouter) LoadTestAccount(c echo.Context) error {
	id := c.Param("id")
	account, err := router.readerRepo.LoadSandboxAccount(id)
	if err != nil {
		return render.NewDBError(err)
	}

	if !account.IsTest() {
		return render.NewNotFound("Not Found")
	}

	return c.JSON(http.StatusOK, account)
}

// DeleteTestAccount deletes a sandbox account.
//
// DELETE /sandbox/:id
func (router ReaderRouter) DeleteTestAccount(c echo.Context) error {
	id := c.Param("id")

	found, err := router.readerRepo.SandboxUserExists(id)
	if err != nil {
		return render.NewDBError(err)
	}

	if !found {
		return c.NoContent(http.StatusNoContent)
	}

	err = router.readerRepo.DeleteTestAccount(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeSandboxPassword overrides current password.
// Input:
// ftcId: string;
// password: string;
func (router ReaderRouter) ChangeSandboxPassword(c echo.Context) error {
	id := c.Param("id")

	var input reader.TestPasswordUpdater
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	input.FtcID = id

	err := router.readerRepo.ChangePassword(input)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
