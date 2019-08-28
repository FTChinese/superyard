package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
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

// LoadFTCAccount retrieves a ftc user's profile.
//
//	GET /readers/ftc/{id}
func (router ReaderRouter) LoadFTCAccount(w http.ResponseWriter, req *http.Request) {
	ftcID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	a, err := router.env.LoadAccountFtc(ftcID)

	// 404 Not Found
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	_ = view.Render(w, view.NewResponse().SetBody(a))
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

	lh, err := router.env.ListLoginHistory(userID, pagination)
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

	a, err := router.env.LoadAccountWx(unionID)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	// 200 OK
	_ = view.Render(w, view.NewResponse().SetBody(a))
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
	ah, err := router.env.ListOAuthHistory(unionID, pagination)

	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(ah))
}

func (router ReaderRouter) GetOrder(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	order, err := router.env.RetrieveOrder(id)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(order))
}
