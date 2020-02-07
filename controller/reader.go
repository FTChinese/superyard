package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/builder"
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/customer"
	"net/http"
)

// ReaderRouter responds to requests for customer services.
type ReaderRouter struct {
	env customer.Env
}

// NewReaderRouter creates a new instance of ReaderRouter
func NewReaderRouter(db *sqlx.DB) ReaderRouter {
	return ReaderRouter{
		env: customer.Env{DB: db},
	}
}

type accountResult struct {
	success reader.Account
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

	account := accountResult.success
	account.Membership = memberResult.success
	account.Kind = reader.AccountKindFtc

	// 200 OK
	return c.JSON(http.StatusOK, account)
}

// LoadLoginHistory retrieves a list of login history.
//
// GET /users/ftc/login-history/{id}?page=<number>&per_page=<number>
func (router ReaderRouter) LoadLoginHistory(c echo.Context) error {

	userID := c.Param("id")

	var pagination builder.Pagination
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

	account := accountResult.success
	account.Membership = memberResult.success
	account.Kind = reader.AccountKindWx

	// 200 OK
	return c.JSON(http.StatusOK, account)
}

// LoadOAuthHistory retrieves a wechat user oauth history.
//
// GET /users/wx/{id}/login?page=<number>&per_page=<number>
func (router ReaderRouter) LoadOAuthHistory(c echo.Context) error {

	unionID := c.Param("id")

	var pagination builder.Pagination
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
