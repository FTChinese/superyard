package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/view"
	"github.com/jmoiron/sqlx"
	"gitlab.com/ftchinese/backyard-api/models/reader"
	"gitlab.com/ftchinese/backyard-api/repository/customer"
	"net/http"
)

type OrderRouter struct {
	env customer.Env
}

func NewOrderRouter(db *sqlx.DB) OrderRouter {
	return OrderRouter{
		env: customer.Env{DB: db},
	}
}

// ListOrders shows a list of a user's orders
func (router OrderRouter) ListOrders(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	q := struct {
		FtcID   string `schema:"ftc_id"`
		UnionID string `schema:"union_id"`
		Page    int64  `schema:"page"`
		PerPage int64  `schema:"per_page"`
	}{}

	if err := decoder.Decode(&q, req.Form); err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	accountID := reader.NewAccountID(q.FtcID, q.UnionID)
	p := gorest.NewPagination(q.Page, q.PerPage)

	orders, err := router.env.ListOrders(accountID, p)
	if err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewResponse().SetBody(orders))
}

func (router OrderRouter) CreateOrder(w http.ResponseWriter, req *http.Request) {
	_, _ = w.Write([]byte("not implemented"))
}

// LoadOrder retrieve an order by id.
func (router OrderRouter) LoadOrder(w http.ResponseWriter, req *http.Request) {
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

// ConfirmOrder set an order confirmation time,
// and create/renew/upgrade membership based on this order.
func (router OrderRouter) ConfirmOrder(w http.ResponseWriter, req *http.Request) {
	orderID, err := GetURLParam(req, "id").ToString()
	if err != nil {
		_ = view.Render(w, view.NewBadRequest(err.Error()))
		return
	}

	if err := router.env.ConfirmOrder(orderID); err != nil {
		_ = view.Render(w, view.NewDBFailure(err))
		return
	}

	_ = view.Render(w, view.NewNoContent())
}
