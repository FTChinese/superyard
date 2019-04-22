package controller

import (
	"database/sql"
	gorest "github.com/FTChinese/go-rest"
	"net/http"

	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/model"

	"github.com/FTChinese/go-rest/view"
)

// UserRouter responds to requests for customer services.
type UserRouter struct {
	model  model.UserEnv
	search model.SearchEnv
}

// NewUserRouter creates a new instance of UserRouter
func NewUserRouter(db *sql.DB) UserRouter {
	return UserRouter{
		search: model.SearchEnv{DB: db},
		model:  model.UserEnv{DB: db},
	}
}

// LoadFTCAccount retrieves a ftc user's profile.
//
//	GET /users/ftc/account/{id}
func (router UserRouter) LoadFTCAccount(w http.ResponseWriter, req *http.Request) {
	userID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	a, err := router.model.LoadAccountByID(userID)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().SetBody(a))
}

// LoadLoginHistory retrieves a list of login history.
//
// GET /users/ftc/login-history/{id}?page=<number>&per_page=<number>
func (router UserRouter) LoadLoginHistory(w http.ResponseWriter, req *http.Request) {

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

	lh, err := router.model.ListLoginHistory(userID, pagination)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(lh))
}

// LoadOrders list all order placed by a user.
//
//	GET /users/ftc/orders/{id}?page=<number>&per_page=<number>
func (router UserRouter) LoadOrders(w http.ResponseWriter, req *http.Request) {

	userID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	u, err := router.search.FindUserByID(userID)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	pagination := gorest.GetPagination(req)

	orders, err := router.model.ListOrders(
		null.StringFrom(u.UserID),
		u.UnionID,
		pagination)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().NoCache().SetBody(orders))
}

// LoadOrdersWxOnly list orders placed by a wechat-only account.
//
//	GET /users/wx/orders/{id}?page=<number>&per_page=<number>
func (router UserRouter) LoadOrdersWxOnly(w http.ResponseWriter, req *http.Request) {
	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	orders, err := router.model.ListOrders(
		null.String{},
		null.StringFrom(unionID),
		pagination)

	// 404 Not Found
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
func (router UserRouter) LoadWxAccount(w http.ResponseWriter, req *http.Request) {
	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		logger.WithField("trace", "LoadWxAccount").Error(err)
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	a, err := router.model.LoadAccountByWx(unionID)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().SetBody(a))
}

// LoadOAuthHistory retrieves a wechat user oauth history.
//
// GET /users/wx/oauth-history/{id}?page=<number>&per_page=<number>
func (router UserRouter) LoadOAuthHistory(w http.ResponseWriter, req *http.Request) {

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
	ah, err := router.model.ListOAuthHistory(unionID, pagination)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(ah))
}
