package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/readers"
	"github.com/FTChinese/superyard/internal/app/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/postman"
	"github.com/FTChinese/superyard/pkg/validator"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// ReaderRouter responds to requests for customer services.
type ReaderRouter struct {
	Repo       readers.Env
	Postman    postman.Postman
	APIClients subsapi.APIClients
	Logger     *zap.Logger
	Version    string
}

// LoadFTCAccount retrieves a ftc user's profile.
//
//	GET /readers/ftc/:id
func (router ReaderRouter) LoadFTCAccount(c echo.Context) error {
	defer router.Logger.Sync()
	sugar := router.Logger.Sugar()

	ftcID := c.Param("id")

	resp, err := router.APIClients.
		Select(true).
		LoadFtcAccount(ftcID)

	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// LoadWxAccount retrieves a wechat user's account
//
//	GET /users/wx/account/:id
func (router ReaderRouter) LoadWxAccount(c echo.Context) error {
	defer router.Logger.Sync()
	sugar := router.Logger.Sugar()

	unionID := c.Param("id")

	resp, err := router.APIClients.
		Select(true).
		LoadWxAccount(unionID)

	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) LoadFtcProfile(c echo.Context) error {

	defer router.Logger.Sync()
	sugar := router.Logger.Sugar()

	ftcID := c.Param("id")

	resp, err := router.APIClients.
		Select(true).
		LoadFtcProfile(ftcID)

	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	if err != nil {
		sugar.Error(err)
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
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
