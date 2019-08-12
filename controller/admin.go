package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/repository"
	"gitlab.com/ftchinese/backyard-api/repository/admin"
	"gitlab.com/ftchinese/backyard-api/repository/staff"
	"net/http"

	"github.com/FTChinese/go-rest/view"
)

// AdminRouter responds to administration tasks performed by a superuser.
type AdminRouter struct {
	env      admin.Env // used by administrator to retrieve staff profile
	search   repository.SearchEnv
	postman  postoffice.Postman
	staffEnv staff.Env
}

// NewAdminRouter creates a new instance of AdminRouter.
func NewAdminRouter(db *sqlx.DB, p postoffice.Postman) AdminRouter {
	return AdminRouter{
		env:      admin.Env{DB: db},
		search:   repository.SearchEnv{DB: db},
		staffEnv: staff.Env{DB: db},
		postman:  p,
	}
}

// FindAccount searches staff account either by name or by email.
//
//	GET /admin/account/search?k={name|email}&v={value}
func (router AdminRouter) SearchStaff(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	key, err := gorest.GetQueryParam(req, "k").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}
	val, err := gorest.GetQueryParam(req, "v").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	var col staff.Column
	switch key {
	case "name":
		col = staff.ColumnUserName
	case "email":
		col = staff.ColumnEmail

	// `400 Bad Request`
	default:
		resp := view.NewBadRequest("The value of 'k' must be one of 'name' or 'email'")
		view.Render(w, resp)
		return
	}

	profile, err := router.staffEnv.Load(col, val)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(profile))
}

// ListVIP lists all ftc account granted vip.
//
//	GET /admin/vip?page=<number>&per_page=<number>
func (router AdminRouter) ListVIP(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()

	// 400 Bad Request if query string cannot be parsed.
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	myfts, err := router.env.ListVIP(pagination)

	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(myfts))
}

// GrantVIP grants vip to an ftc account.
//
//	PUT /admin/vip/{id}
func (router AdminRouter) GrantVIP(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err = router.env.GrantVIP(id); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}

// RevokeVIP removes a ftc account from vip.
//
//	DELETE /admin/vip/{id}
func (router AdminRouter) RevokeVIP(w http.ResponseWriter, req *http.Request) {
	id, err := GetURLParam(req, "id").ToString()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := router.env.RevokeVIP(id); err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	// 204 No Content
	view.Render(w, view.NewNoContent())
}
