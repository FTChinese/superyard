package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/reader"
	"gitlab.com/ftchinese/backyard-api/repository/search"
	"net/http"
)

type SearchRouter struct {
	env search.Env
}

func NewSearchRouter(db *sqlx.DB) SearchRouter {
	return SearchRouter{
		env: search.Env{DB: db},
	}
}

// SearchFtcUser tries to find a user by userName or email
//
//	GET /search/reader?email=<name@example.org>
func (router SearchRouter) SearchFtcUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	// 400 Bad Request
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	var param reader.SearchParam
	if err := decoder.Decode(&param, req.Form); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}
	param.Sanitize()
	if err := param.RequireEmail(); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	ftcInfo, err := router.env.SearchFtcUser(param.Email)

	// 404 Not Found
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(ftcInfo))
}

// FindWxUser tries to find a wechat user by nickname
//
// GET /search/reader/wx?q=<nickname>&page=<number>&per_page=<number>
func (router SearchRouter) SearchWxUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	var param reader.SearchParam
	if err := decoder.Decode(&param, req.Form); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}
	param.Sanitize()
	if err := param.RequireQ(); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	pagination := gorest.GetPagination(req)

	wxUsers, err := router.env.SearchWxUser(param.Q, pagination)

	// 404 Not Found
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().NoCache().SetBody(wxUsers))
}
