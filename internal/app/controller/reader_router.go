package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/readers"
	"github.com/FTChinese/superyard/internal/app/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/postman"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// ReaderRouter responds to requests for customer services.
type ReaderRouter struct {
	Repo       readers.Env
	Postman    postman.Postman
	APIClient  subsapi.Client // Deprecated
	APIClients subsapi.APIClients
	Logger     *zap.Logger
	Version    string
}

// FindFTCAccount searches an ftc account by email or user name.
//
// GET /readers/ftc?q=<email|username>
func (router ReaderRouter) FindFTCAccount(c echo.Context) error {
	value := strings.TrimSpace(c.QueryParam("q"))
	if value == "" {
		return render.NewBadRequest("Missing query parameter q")
	}

	a, err := router.Repo.FindFtcAccount(value)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, a)
}

// LoadFTCAccount retrieves a ftc user's profile.
//
//	GET /readers/ftc/:id
func (router ReaderRouter) LoadFTCAccount(c echo.Context) error {
	defer router.Logger.Sync()
	sugar := router.Logger.Sugar()

	ftcID := c.Param("id")

	account, err := router.Repo.AccountByFtcID(ftcID)

	if err != nil {
		sugar.Error(err)
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// LoadActivities retrieves a list of login history.
//
// GET /reader/ftc/:id/activities?page=<number>&per_page=<number>
func (router ReaderRouter) LoadActivities(c echo.Context) error {

	ftcID := c.Param("id")

	var pagination gorest.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	lh, err := router.Repo.ListActivities(ftcID, pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, lh)
}

// LoadWxAccount retrieves a wechat user's account
//
//	GET /users/wx/account/:id
func (router ReaderRouter) LoadWxAccount(c echo.Context) error {
	unionID := c.Param("id")

	account, err := router.Repo.AccountByUnionID(unionID)
	if err != nil {
		return render.NewDBError(err)
	}

	if !account.IsTest() {
		return render.NewNotFound("Not Found")
	}

	return c.JSON(http.StatusOK, account)
}

// LoadOAuthHistory retrieves a wechat user oauth history.
//
// GET /users/wx/:id/login?page=<number>&per_page=<number>
func (router ReaderRouter) LoadOAuthHistory(c echo.Context) error {

	unionID := c.Param("id")

	var pagination gorest.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	ah, err := router.Repo.ListWxLoginHistory(unionID, pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, ah)
}

func (router ReaderRouter) LoadFtcProfile(c echo.Context) error {
	ftcID := c.Param("id")

	p, err := router.Repo.RetrieveFtcProfile(ftcID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, p)
}

func (router ReaderRouter) LoadWxProfile(c echo.Context) error {
	unionID := c.Param("id")

	p, err := router.Repo.RetrieveWxProfile(unionID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, p)
}

// SearchAccount tries to find a reader's account.
// Query parameters: q=<email | nickname>&kind=<ftc | wechat>&page=<number>&per_page=<number>
func (router ReaderRouter) SearchAccount(c echo.Context) error {
	q := c.QueryParam("q")
	k := c.QueryParam("kind")
	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}
	page.Normalize()

	switch k {
	case "ftc":
		if ve := validator.New("q").Required().Email().Validate(q); ve != nil {
			return render.NewUnprocessable(ve)
		}

		accounts, err := router.Repo.SearchJoinedAccountEmail(q, page)
		if err != nil {
			return render.NewDBError(err)
		}
		// Email is always uniquely constrained, therefore at most one item is retrieved.
		return c.JSON(http.StatusOK, accounts)

	case "wechat":

		accounts, err := router.Repo.SearchJoinedAccountWxName(q, page)
		if err != nil {
			return render.NewDBError(err)
		}

		return c.JSON(http.StatusOK, accounts)

	default:
		return render.NewBadRequest("Query account kind could only be ftc or wechat")
	}
}
