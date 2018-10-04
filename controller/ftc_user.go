package controller

import (
	"net/http"

	"gitlab.com/ftchinese/backyard-api/ftcuser"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// FTCUserController handles various customer service tasks
type FTCUserController struct {
	model ftcuser.Env
}

// SearchUser tries to find a user by userName or email
// /search/user?k=<name|email>&v=:value
func (c FTCUserController) SearchUser(w http.ResponseWriter, req *http.Request) {
	key := req.Form.Get("k")
	val := req.Form.Get("v")

	if key == "" || val == "" {
		resp := util.NewBadRequest("Both 'k' and 'v' should be present in query string")
		view.Render(w, resp)

		return
	}

	var a ftcuser.Account
	var err error
	switch key {
	case "name":
		a, err = c.model.FindUserByName(val)

	case "email":
		a, err = c.model.FindUserByEmail(val)

	default:
		resp := util.NewBadRequest("The value of 'k' must be one of 'name' or 'email'")
		view.Render(w, resp)
		return
	}

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(a))
}

// UserProfile retrieves a user profile
func (c FTCUserController) UserProfile(w http.ResponseWriter, req *http.Request) {
	userID := getURLParam(req, "userId").toString()

	// 400 Bad Request
	// { message: "Invalid request URI" }
	if userID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	p, err := c.model.Profile(userID)

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(p))
}

// UserOrders list all order placed by a user
func (c FTCUserController) UserOrders(w http.ResponseWriter, req *http.Request) {
	userID := getURLParam(req, "userId").toString()

	// 400 Bad Request
	// { message: "Invalid request URI" }
	if userID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	orders, err := c.model.OrderRoster(userID)

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(orders))
}

// LoginHistory lists a user's login footprint
func (c FTCUserController) LoginHistory(w http.ResponseWriter, req *http.Request) {
	userID := getURLParam(req, "userId").toString()

	// 400 Bad Request
	// { message: "Invalid request URI" }
	if userID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	history, err := c.model.LoginHistory(userID)

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err, ""))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(history))
}
