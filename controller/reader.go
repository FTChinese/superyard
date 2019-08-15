package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/reader"
	"gitlab.com/ftchinese/backyard-api/repository/customer"
	"net/http"

	"github.com/FTChinese/go-rest/view"
	"github.com/guregu/null"
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

// LoadFTCAccount retrieves a ftc user's profile.
//
//	GET /readers/ftc/{id}
func (router ReaderRouter) LoadFTCAccount(w http.ResponseWriter, req *http.Request) {
	ftcID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	a, err := router.env.LoadAccountFtc(ftcID)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().SetBody(a))
}

// LoadOrders list all order placed by a user.
//
//	GET /readers/ftc/{id}/orders?page=<number>&per_page=<number>
func (router ReaderRouter) LoadFtcOrders(w http.ResponseWriter, req *http.Request) {

	ftcID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	orders, err := router.env.ListOrders(
		reader.AccountID{
			CompoundID: ftcID,
			FtcID:      null.StringFrom(ftcID),
			UnionID:    null.String{},
		},
		pagination)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().NoCache().SetBody(orders))
}

// LoadLoginHistory retrieves a list of login history.
//
// GET /users/ftc/login-history/{id}?page=<number>&per_page=<number>
func (router ReaderRouter) LoadLoginHistory(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	userID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	lh, err := router.env.ListLoginHistory(userID, pagination)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(lh))
}

// LoadOrdersWxOnly list orders placed by a wechat-only account.
//
//	GET /users/wx/{id}/orders/?page=<number>&per_page=<number>
func (router ReaderRouter) LoadWxOrders(w http.ResponseWriter, req *http.Request) {
	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	orders, err := router.env.ListOrders(
		reader.AccountID{
			CompoundID: unionID,
			FtcID:      null.String{},
			UnionID:    null.StringFrom(unionID),
		},
		pagination)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().NoCache().SetBody(orders))
}

// LoadWxAccount retrieves a wechat user's account
//
//	GET /users/wx/account/{id}
func (router ReaderRouter) LoadWxAccount(w http.ResponseWriter, req *http.Request) {
	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		logger.WithField("trace", "ReaderRouter.LoadWxAccount").Error(err)
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	a, err := router.env.LoadAccountWx(unionID)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().SetBody(a))
}

// LoadOAuthHistory retrieves a wechat user oauth history.
//
// GET /users/wx/{id}/login?page=<number>&per_page=<number>
func (router ReaderRouter) LoadOAuthHistory(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)
	ah, err := router.env.ListOAuthHistory(unionID, pagination)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(ah))
}

func (router ReaderRouter) GetOrder(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	order, err := router.env.RetrieveOrder(id)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(order))
}
