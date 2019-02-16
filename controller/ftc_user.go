package controller

import (
	"database/sql"
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

// LoadAccount retrieves a ftc user's profile.
//
//	GET /users/{id}/account
func (router UserRouter) LoadAccount(w http.ResponseWriter, req *http.Request) {
	userID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p, err := router.model.LoadAccountByID(userID)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().SetBody(p))
}

// LoadLoginHistory retrieves a list of login history.
// GET /users/{id}/login-history
func (router UserRouter) LoadLoginHistory(w http.ResponseWriter, req *http.Request) {
	userID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	lh, err := router.model.ListLoginHistory(userID)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(lh))
}

// LoadOrders list all order placed by a user.
//
//	GET /users/{id}/orders
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

	orders, err := router.model.ListOrders(null.StringFrom(u.UserID), u.UnionID)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().NoCache().SetBody(orders))
}

// LoadWxInfo retrieves a ftc user's profile.
//
//	GET /wxusers/{id}/info
func (router UserRouter) LoadWxInfo(w http.ResponseWriter, req *http.Request) {
	unionID, err := GetQueryParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	p, err := router.model.LoadWxInfo(unionID)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().SetBody(p))
}

// LoadOAuthHistory retrieves a wechat user oauth history.
//
// GET /wxusers/{id}/oauth-history
func (router UserRouter) LoadOAuthHistory(w http.ResponseWriter, req *http.Request) {
	unionID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	ah, err := router.model.ListOAuthHistory(unionID)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(ah))
}
