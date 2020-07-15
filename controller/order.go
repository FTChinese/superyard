package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/pkg/subs"
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

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	var ids subs.CompoundIDs
	if err := c.Bind(&ids); err != nil {
		return render.NewBadRequest(err.Error())
	}

	orders, err := router.env.ListOrders(ids, page)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, orders)
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

	err := router.env.ConfirmOrder(orderID)

	if err != nil {
		switch err {
		case subs.ErrAlreadyConfirmed:
			// Order already confirmed.
			return render.NewUnprocessable(&render.ValidationError{
				Message: err.Error(),
				Field:   "confirmedAt",
				Code:    render.CodeAlreadyExists,
			})

		case subs.ErrValidNonAliOrWxPay:
			// A valid membership not purchased via FTC order.
			return render.NewUnprocessable(&render.ValidationError{
				Message: err.Error(),
				Field:   "membership",
				Code:    "non_expired_non_ftc",
			})

		case subs.ErrAlreadyUpgraded:
			// Membership is already a premium.
			return render.NewUnprocessable(&render.ValidationError{
				Message: err.Error(),
				Field:   "membership",
				Code:    "already_premium",
			})

		default:
			return render.NewDBError(err)
		}
	}

	return c.NoContent(http.StatusNoContent)
}
