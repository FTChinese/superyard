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

type accountResult struct {
	success reader.BaseAccount
	err     error
}

type memberResult struct {
	success reader.Membership
	err     error
}

// LoadFTCAccount retrieves a ftc user's profile.
//
//	GET /readers/ftc/{id}
func (router ReaderRouter) LoadFTCAccount(c echo.Context) error {
	ftcID := c.Param("id")

	aChan := make(chan accountResult)
	mChan := make(chan memberResult)

	go func() {
		account, err := router.env.RetrieveAccountFtc(ftcID)
		aChan <- accountResult{
			success: account,
			err:     err,
		}
	}()

	go func() {
		member, err := router.env.RetrieveMemberFtc(ftcID)
		mChan <- memberResult{
			success: member,
			err:     err,
		}
	}()

	accountResult, memberResult := <-aChan, <-mChan
	if accountResult.err != nil {
		return util.NewDBFailure(accountResult.err)
	}
	if memberResult.err != nil {
		return util.NewDBFailure(memberResult.err)
	}

	// 200 OK
	return c.JSON(http.StatusOK, reader.Account{
		BaseAccount: accountResult.success,
		Membership:  memberResult.success,
	})
}

// LoadLoginHistory retrieves a list of login history.
//
// GET /users/ftc/login-history/{id}?page=<number>&per_page=<number>
func (router ReaderRouter) LoadLoginHistory(c echo.Context) error {

	userID := c.Param("id")

	var pagination util.Pagination
	if err := c.Bind(&pagination); err != nil {
		return util.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	lh, err := router.env.ListEmailLoginHistory(userID, pagination)
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

	aChan := make(chan accountResult)
	mChan := make(chan memberResult)

	go func() {
		a, err := router.env.RetrieveAccountWx(unionID)
		aChan <- accountResult{
			success: a,
			err:     err,
		}
	}()

	go func() {
		m, err := router.env.RetrieveMemberWx(unionID)
		mChan <- memberResult{
			success: m,
			err:     err,
		}
	}()

	accountResult, memberResult := <-aChan, <-mChan
	if accountResult.err != nil {
		return util.NewDBFailure(accountResult.err)
	}
	if memberResult.err != nil {
		return util.NewDBFailure(memberResult.err)
	}

	// 200 OK
	return c.JSON(http.StatusOK, reader.Account{
		BaseAccount: accountResult.success,
		Membership:  memberResult.success,
	})
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
		// Find ftc id by email
		ftcID, err := router.env.SearchFtcIDByEmail(q)
		// The email might not exist.
		if err != nil {
			return util.NewDBFailure(err)
		}
		// Email is always uniquely constrained, therefore at most one item is retrieved.
		a, err := router.env.RetrieveAccountFtc(ftcID)
		return c.JSON(http.StatusOK, []reader.BaseAccount{a})

	case "wechat":
		var p util.Pagination
		if err := c.Bind(&p); err != nil {
			return util.NewBadRequest(err.Error())
		}
		p.Normalize()

		unionIDs, err := router.env.SearchWxIDs(q, p)
		if err != nil {
			return util.NewDBFailure(err)
		}

		accounts, err := router.env.RetrieveWxAccounts(unionIDs)
		if err != nil {
			return util.NewDBFailure(err)
		}

		return c.JSON(http.StatusOK, accounts)

	default:
		return util.NewBadRequest("Query account kind could only be ftc or wechat")
	}
}
