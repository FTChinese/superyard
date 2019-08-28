package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/reader"
	"gitlab.com/ftchinese/backyard-api/repository/customer"
	"net/http"

	"github.com/FTChinese/go-rest/view"
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
func (router ReaderRouter) LoadFTCAccount(w http.ResponseWriter, req *http.Request) {
	ftcID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

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
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}
	if memberResult.err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	account := accountResult.success
	account.Membership = memberResult.success

	// 200 OK
	_ = view.Render(w, view.NewResponse().SetBody(account))
}

// LoadLoginHistory retrieves a list of login history.
//
// GET /users/ftc/login-history/{id}?page=<number>&per_page=<number>
func (router ReaderRouter) LoadLoginHistory(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	userID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	lh, err := router.env.ListEmailLoginHistory(userID, pagination)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(lh))
}

// LoadWxAccount retrieves a wechat user's account
//
//	GET /users/wx/account/{id}
func (router ReaderRouter) LoadWxAccount(w http.ResponseWriter, req *http.Request) {
	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		logger.WithField("trace", "ReaderRouter.LoadWxAccount").Error(err)
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

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
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}
	if memberResult.err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	account := accountResult.success
	account.Membership = memberResult.success

	// 200 OK
	_ = view.Render(w, view.NewResponse().SetBody(account))
}

// LoadOAuthHistory retrieves a wechat user oauth history.
//
// GET /users/wx/{id}/login?page=<number>&per_page=<number>
func (router ReaderRouter) LoadOAuthHistory(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)
	ah, err := router.env.ListWxLoginHistory(unionID, pagination)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(ah))
}

func (router ReaderRouter) LoadFtcProfile(w http.ResponseWriter, req *http.Request) {
	ftcID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p, err := router.env.RetrieveFtcProfile(ftcID)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(p))
}

func (router ReaderRouter) LoadWxProfile(w http.ResponseWriter, req *http.Request) {
	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		logger.WithField("trace", "ReaderRouter.LoadWxProfile").Error(err)
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p, err := router.env.RetrieveWxProfile(unionID)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(p))
}
