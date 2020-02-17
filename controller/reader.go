package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/models/validator"
	"gitlab.com/ftchinese/superyard/repository/readers"
	"net/http"
)

// ReaderRouter responds to requests for customer services.
type ReaderRouter struct {
	env readers.Env
}

// NewReaderRouter creates a new instance of ReaderRouter
func NewReaderRouter(db *sqlx.DB) ReaderRouter {
	return ReaderRouter{
		env: readers.Env{DB: db},
	}
}

// LoadFTCAccount retrieves a ftc user's profile.
//
//	GET /readers/ftc/{id}
func (router ReaderRouter) LoadFTCAccount(c echo.Context) error {
	ftcID := c.Param("id")

	account, err := router.env.LoadFTCAccount(ftcID)

	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, account)
}

// LoadActivities retrieves a list of login history.
//
// GET /users/ftc/login-history/{id}?page=<number>&per_page=<number>
func (router ReaderRouter) LoadActivities(c echo.Context) error {

	ftcID := c.Param("id")

	var pagination util.Pagination
	if err := c.Bind(&pagination); err != nil {
		return util.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	lh, err := router.env.ListActivities(ftcID, pagination)
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, lh)
}

// LoadWxAccount retrieves a wechat user's account
//
//	GET /users/wx/account/{id}
func (router ReaderRouter) LoadWxAccount(c echo.Context) error {
	unionID := c.Param("id")

	account, err := router.env.LoadWxAccount(unionID)
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, account)
}

// LoadOAuthHistory retrieves a wechat user oauth history.
//
// GET /users/wx/{id}/login?page=<number>&per_page=<number>
func (router ReaderRouter) LoadOAuthHistory(c echo.Context) error {

	unionID := c.Param("id")

	var pagination util.Pagination
	if err := c.Bind(&pagination); err != nil {
		return util.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	ah, err := router.env.ListWxLoginHistory(unionID, pagination)
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, ah)
}

func (router ReaderRouter) LoadFtcProfile(c echo.Context) error {
	ftcID := c.Param("id")

	p, err := router.env.RetrieveFtcProfile(ftcID)
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, p)
}

func (router ReaderRouter) LoadWxProfile(c echo.Context) error {
	unionID := c.Param("id")

	p, err := router.env.RetrieveWxProfile(unionID)
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, p)
}

// SearchAccount tries to find a reader's account.
// Query parameters: q=<email | nickname>&kind=<ftc | wechat>&page=<number>&per_page=<number>
func (router ReaderRouter) SearchAccount(c echo.Context) error {
	q := c.QueryParam("q")
	k := c.QueryParam("kind")

	switch k {
	case "ftc":
		if ie := validator.New("q").Required().Email().Validate(q); ie != nil {
			return util.NewUnprocessable(ie)
		}

		a, err := router.env.SearchFtcAccount(q)
		if err != nil {
			return util.NewDBFailure(err)
		}
		// Email is always uniquely constrained, therefore at most one item is retrieved.
		return c.JSON(http.StatusOK, []reader.BaseAccount{a})

	case "wechat":
		var p util.Pagination
		if err := c.Bind(&p); err != nil {
			return util.NewBadRequest(err.Error())
		}
		p.Normalize()

		accounts, err := router.env.SearchWxAccounts(q, p)
		if err != nil {
			return util.NewDBFailure(err)
		}

		return c.JSON(http.StatusOK, accounts)

	default:
		return util.NewBadRequest("Query account kind could only be ftc or wechat")
	}
}
