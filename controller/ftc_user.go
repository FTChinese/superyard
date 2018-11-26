package controller

import (
	"database/sql"
	"net/http"

	"gitlab.com/ftchinese/backyard-api/ftcuser"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

// FTCUserRouter responds to requests for customer services.
type FTCUserRouter struct {
	model ftcuser.Env
}

// NewFTCUserRouter creates a new instance of FTCUserRouter
func NewFTCUserRouter(db *sql.DB) FTCUserRouter {
	model := ftcuser.Env{DB: db}

	return FTCUserRouter{
		model: model,
	}
}

// SearchUser tries to find a user by userName or email
//
//	GET /search/user?k=<name|email>&v=<value>
func (c FTCUserRouter) SearchUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	// 400 Bad Request
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))

		return
	}

	key := req.Form.Get("k")
	val := req.Form.Get("v")

	if key == "" || val == "" {
		resp := util.NewBadRequest("Both 'k' and 'v' should be present in query string")
		view.Render(w, resp)

		return
	}

	var a ftcuser.Account
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
		view.Render(w, util.NewDBFailure(err))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(a))
}

// UserProfile retrieves a ftc user's profile.
//
//	GET /ftc-user/profile/{userId}
func (c FTCUserRouter) UserProfile(w http.ResponseWriter, req *http.Request) {
	userID := getURLParam(req, "userId").toString()

	// 400 Bad Request
	if userID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	p, err := c.model.Profile(userID)

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, util.NewResponse().NoCache().SetBody(p))
}

// UserOrders list all order placed by a user.
//
//	GET /ftc-user/profile/{userId}/orders
//
// - `400 Bad Request` if request URL does not contain `userId` part
//	{
//		"message": "Invalid request URI"
//	}
//
// - 200 OK with body:
// 	[{
// 		"orderId": "",
// 		"tierToBuy": "standard | premium",
// 		"price": 198.00,
// 		"totalAmount": 198.00,
// 		"billingCycle": "year | month",
// 		"paymentMethod": "alipay | tenpay | stripe | redeem_code",
// 		"clientType": "web | ios | android | unknown",
// 		"clientVersion": "1.2.1",
// 		"createdAt": "",
// 		"confirmedAt": "",
// 		"userIp": "127.0.0.1"
// 	}]
func (c FTCUserRouter) UserOrders(w http.ResponseWriter, req *http.Request) {
	userID := getURLParam(req, "userId").toString()

	// 400 Bad Request
	if userID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	orders, err := c.model.OrderRoster(userID)

	// 404 Not Found
	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	// 200 OK
	view.Render(w, util.NewResponse().NoCache().SetBody(orders))
}

// LoginHistory lists a user's login history. 100 entries per page.
//
//	GET /ftc-user/profile/{userId}/login?page=<number>
//
// - `400 Bad Request` if query string cannot be parsed, or if request URL does not contain `userId` part
//	{
//		"message": "Invalid request URI"
//	}
//
// - 200 OK with body:
// [{
// 		"authMethod": "email | phone | wechat | weibo",
// 		"clientType": "web | ios | android | unknown",
// 		"clientVersion": "3.1.2",
// 		"userIp": "127.0.0.1",
// 		"loggedInAt": ""
// }]
func (c FTCUserRouter) LoginHistory(w http.ResponseWriter, req *http.Request) {
	// Parse form
	err := req.ParseForm()

	// 400 Bad Request if query string cannot be parsed.
	if err != nil {
		view.Render(w, util.NewBadRequest(err.Error()))
		return
	}

	// Get page query
	page, err := getQueryParam(req, "page").toInt()

	if err != nil {
		page = 1
	}

	// Get userId
	userID := getURLParam(req, "userId").toString()

	// 400 Bad Request
	if userID == "" {
		view.Render(w, util.NewBadRequest("Invalid request URI"))

		return
	}

	history, err := c.model.LoginHistory(userID, page, 100)

	if err != nil {
		view.Render(w, util.NewDBFailure(err))
		return
	}

	view.Render(w, util.NewResponse().NoCache().SetBody(history))
}
