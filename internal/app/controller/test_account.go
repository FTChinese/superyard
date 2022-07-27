package controller

import (
	"encoding/json"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/pkg/sandbox"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateTestUser creates a test account.
//
// POST /sandbox
//
// Input: reader.TestAccountInput
// email: string
// password: string
func (router ReaderRouter) CreateTestUser(c echo.Context) error {
	claims := getPassportClaims(c)

	var params sandbox.SignUpParams
	if err := c.Bind(&params); err != nil {
		return render.NewBadRequest(err.Error())
	}
	if ve := params.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	header := xhttp.NewHeaderBuilder().
		WithPlatformWeb().
		WithClientVersion(router.Version).
		WithUserIP(c.RealIP()).
		WithUserAgent(c.Request().UserAgent()).
		Build()

	resp, err := router.APIClients.Select(true).SignUp(params, header)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	if resp.StatusCode != 200 {
		return c.Blob(resp.StatusCode, fetch.ContentJSON, resp.Body)
	}

	var ba sandbox.BaseAccount
	if err := json.Unmarshal(resp.Body, &ba); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ba.FtcID == "" {
		return render.NewInternalError("failed to create account")
	}

	ta := ba.NewTestAccount(params, claims.Username)

	err = router.Repo.CreateTestUser(ta)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, ta)
}

// ListTestUsers retrieves all sandbox user.
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

// LoadTestAccount loads a sandbox account.
//
// GET /sandbox/:id
func (router ReaderRouter) LoadTestAccount(c echo.Context) error {
	id := c.Param("id")
	account, err := router.Repo.LoadSandboxAccount(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// DeleteTestAccount deletes a sandbox account.
//
// DELETE /sandbox/:id
func (router ReaderRouter) DeleteTestAccount(c echo.Context) error {
	id := c.Param("id")

	ta, err := router.Repo.LoadSandboxAccount(id)
	if err != nil {
		return render.NewDBError(err)
	}

	resp, err := router.APIClients.Select(true).DeleteReader(ta)

	if err != nil {
		return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
	}

	err = router.Repo.DeleteTestAccount(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ChangeSandboxPassword overrides current password.
// Input:
// password: string;
func (router ReaderRouter) ChangeSandboxPassword(c echo.Context) error {
	id := c.Param("id")

	var input sandbox.TestAccount
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidatePassword(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	input.FtcID = id

	err := router.Repo.ChangePassword(input)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
