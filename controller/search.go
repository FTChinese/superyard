package controller

import (
	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/model"
	"gitlab.com/ftchinese/backyard-api/user"
	"net/http"
)

type SearchRouter struct {
	model model.SearchEnv
}

// SearchUser tries to find a user by userName or email
//
//	GET /search/user?k=<name|email>&v=<value>
func (router SearchRouter) SearchUser(w http.ResponseWriter, req *http.Request) {
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