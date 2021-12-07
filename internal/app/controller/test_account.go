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

	err := router.Repo.CreateTestUser(account)
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
	p.Normalize()

	users, err := router.Repo.ListTestFtcAccount(p)
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
	account, err := router.Repo.LoadSandboxAccount(id)
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

	found, err := router.Repo.SandboxUserExists(id)
	if err != nil {
		return render.NewDBError(err)
	}

	if !found {
		return c.NoContent(http.StatusNoContent)
	}

	err = router.Repo.DeleteTestAccount(id)
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

	err := router.Repo.ChangePassword(input)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
