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

// LoadAccount retrieves a ftc user's profile. Header `X-User-Id`
//
//	GET /user/account
func (router UserRouter) LoadAccount(w http.ResponseWriter, req *http.Request) {
	email := req.Header.Get(userEmailKey)

	p, err := router.model.LoadAccount(email)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().SetBody(p))
}

// LoadOrders list all order placed by a user. Header `X-User-Id` or `X-UnionId` or both.
//
//	GET /user/orders
func (router UserRouter) LoadOrders(w http.ResponseWriter, req *http.Request) {
	email := req.Header.Get(userEmailKey)

	u, err := router.search.FindUserByEmail(email)
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
