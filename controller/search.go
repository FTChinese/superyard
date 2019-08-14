package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
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

// SearchFTCUser tries to find a user by userName or email
//
//	GET /search/user?k=email&v=<value>
func (router SearchRouter) SearchFTCUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	param := NewSearchParam(req)
	if err := param.NotEmpty(); err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	ftcInfo, err := router.env.SearchFtcUser(param.Value)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(ftcInfo))
}

// FindWxUser tries to find a wechat user by nickname\
//
// GET /search/user/wx?q=<nickname>&page=<number>&per_page=<number>
func (router SearchRouter) SearchWxUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	nickname := req.Form.Get("q")

	if nickname == "" {
		resp := view.NewBadRequest("'q' should should have a value")
		view.Render(w, resp)

		return
	}

	pagination := gorest.GetPagination(req)

	wxUsers, err := router.env.SearchWxUser(nickname, pagination)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(wxUsers))
}
