package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/pkg/validator"
	"gitlab.com/ftchinese/superyard/repository/readers"
	"net/http"
)

// ReaderRouter responds to requests for customer services.
type ReaderRouter struct {
	env     readers.Env
	postman postoffice.PostOffice
}

// NewReaderRouter creates a new instance of ReaderRouter
func NewReaderRouter(db *sqlx.DB, p postoffice.PostOffice) ReaderRouter {
	return ReaderRouter{
		env: readers.Env{
			DB: db,
		},
		postman: p,
	}
}

// LoadFTCAccount retrieves a ftc user's profile.
//
//	GET /readers/ftc/:id
func (router ReaderRouter) LoadFTCAccount(c echo.Context) error {
	ftcID := c.Param("id")

	account, err := router.env.LoadFTCAccount(ftcID)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// LoadActivities retrieves a list of login history.
//
// GET /reader/ftc//:id/activities?page=<number>&per_page=<number>
func (router ReaderRouter) LoadActivities(c echo.Context) error {

	ftcID := c.Param("id")

	var pagination gorest.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	lh, err := router.env.ListActivities(ftcID, pagination)
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

	account, err := router.env.LoadWxAccount(unionID)
	if err != nil {
		return render.NewDBError(err)
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

	ah, err := router.env.ListWxLoginHistory(unionID, pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, ah)
}

func (router ReaderRouter) LoadFtcProfile(c echo.Context) error {
	ftcID := c.Param("id")

	p, err := router.env.RetrieveFtcProfile(ftcID)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, p)
}

func (router ReaderRouter) LoadWxProfile(c echo.Context) error {
	unionID := c.Param("id")

	p, err := router.env.RetrieveWxProfile(unionID)
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

		accounts, err := router.env.SearchFtcAccount(q, page)
		if err != nil {
			return render.NewDBError(err)
		}
		// Email is always uniquely constrained, therefore at most one item is retrieved.
		return c.JSON(http.StatusOK, accounts)

	case "wechat":

		accounts, err := router.env.SearchWxAccounts(q, page)
		if err != nil {
			return render.NewDBError(err)
		}

		return c.JSON(http.StatusOK, accounts)

	default:
		return render.NewBadRequest("Query account kind could only be ftc or wechat")
	}
}
