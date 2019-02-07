package controller

import (
	"database/sql"
	"github.com/guregu/null"
	"gitlab.com/ftchinese/backyard-api/model"
	"net/http"

	"github.com/FTChinese/go-rest/view"
)

// UserRouter responds to requests for customer services.
type UserRouter struct {
	model model.UserEnv
}

// NewUserRouter creates a new instance of UserRouter
func NewUserRouter(db *sql.DB) UserRouter {
	return UserRouter{
		model: model.UserEnv{DB: db},
	}
}

// LoadAccount retrieves a ftc user's profile. Header `X-User-Id`
//
//	GET /user/account
func (router UserRouter) LoadAccount(w http.ResponseWriter, req *http.Request) {
	userID := req.Header.Get(userIDKey)

	p, err := router.model.LoadAccount(userID)

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
	uID := req.Header.Get(userIDKey)
	wID := req.Header.Get(unionIDKey)

	userID := null.NewString(uID, uID != "")
	unionID := null.NewString(wID, wID != "")

	orders, err := router.model.ListOrders(userID, unionID)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, view.NewResponse().NoCache().SetBody(orders))
}


