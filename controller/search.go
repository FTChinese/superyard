package controller

import (
	"database/sql"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/model"
	"gitlab.com/ftchinese/backyard-api/user"
	"net/http"
)

type SearchRouter struct {
	model model.SearchEnv
}

func NewSearchRouter(db *sql.DB) SearchRouter {
	return SearchRouter{
		model: model.SearchEnv{DB: db},
	}
}

// SearchFTCUser tries to find a user by userName or email
//
//	GET /search/user?k=<name|email>&v=<value>
func (router SearchRouter) SearchFTCUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))

		return
	}

	key := req.Form.Get("k")
	val := req.Form.Get("v")

	if key == "" || val == "" {
		resp := view.NewBadRequest("Both 'k' and 'v' should be present in query string")
		view.Render(w, resp)

		return
	}

	var u user.User
	switch key {
	case "name":
		u, err = router.model.FindUserByName(val)

	case "email":
		u, err = router.model.FindUserByEmail(val)

	default:
		resp := view.NewBadRequest("The value of 'k' must be one of 'name' or 'email'")
		view.Render(w, resp)
		return
	}

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(u))
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

	wxUsers, err := router.model.FindWechat(nickname, pagination)

	// 404 Not Found
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().NoCache().SetBody(wxUsers))
}

// SearchOrder tries to find an order by id.
//
//	GET /search/order?id=<string>
func (router SearchRouter) SearchOrder(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	id, err := gorest.GetQueryParam(req, "id").ToString()
	if err != nil {
		view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	o, err := router.model.FindOrder(id)
	if err != nil {
		view.Render(w, view.NewDBFailure(err))
		return
	}

	view.Render(w, view.NewResponse().SetBody(o))
}
