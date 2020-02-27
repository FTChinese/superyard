package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/repository/readers"
	"net/http"
)

type OrderRouter struct {
	env readers.Env
}

func NewOrderRouter(db *sqlx.DB) OrderRouter {
	return OrderRouter{
		env: readers.Env{DB: db},
	}
}

// ListOrders shows a list of a user's orders
func (router OrderRouter) ListOrders(c echo.Context) error {

	q := struct {
		FtcID   string `query:"ftc_id"`
		UnionID string `query:"union_id"`
		Page    int64  `query:"page"`
		PerPage int64  `query:"per_page"`
	}{}

	if err := c.Bind(&q); err != nil {
		return render.NewBadRequest(err.Error())
	}

	accountID := reader.NewAccountID(q.FtcID, q.UnionID)
	p := gorest.NewPagination(q.Page, q.PerPage)

	orders, err := router.env.ListOrders(accountID, p)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, orders)
}

func (router OrderRouter) CreateOrder(c echo.Context) error {
	return c.String(http.StatusOK, "not implemented")
}

// LoadOrder retrieve an order by id.
func (router OrderRouter) LoadOrder(c echo.Context) error {
	id := c.Param("id")

	order, err := router.env.RetrieveOrder(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, order)
}

// ConfirmOrder set an order confirmation time,
// and create/renew/upgrade membership based on this order.
func (router OrderRouter) ConfirmOrder(c echo.Context) error {
	orderID := c.Param("id")

	if err := router.env.ConfirmOrder(orderID); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
