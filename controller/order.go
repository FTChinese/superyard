package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ListOrders shows a list of a user's orders
// A reader's account might have ftc id or union id, or both.
// We are not sure the account status when the reader creates
// and order, thus client should provide as much ids as possible.
func (router ReaderRouter) ListOrders(c echo.Context) error {

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	var ids subs.CompoundIDs
	if err := c.Bind(&ids); err != nil {
		return render.NewBadRequest(err.Error())
	}

	orders, err := router.readerRepo.ListOrders(ids, page)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, orders)
}

// LoadOrder retrieve an order by id.
func (router ReaderRouter) LoadOrder(c echo.Context) error {
	id := c.Param("id")

	order, err := router.readerRepo.RetrieveOrder(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, order)
}

// ConfirmOrder checks if an order is confirmed, and request API to check its transaction status
// against wxpay or alipay api. After confirmation,
// the corresponding membership is updated by API.
//
// PATCH /orders/:id
func (router ReaderRouter) ConfirmOrder(c echo.Context) error {
	orderID := c.Param("id")

	order, err := router.readerRepo.RetrieveOrder(orderID)
	if err != nil {
		return render.NewDBError(err)
	}

	if order.IsConfirmed() {
		return render.NewUnprocessable(&render.ValidationError{
			Message: "Duplicate confirmation",
			Field:   "confirmedAt",
			Code:    render.CodeAlreadyExists,
		})
	}

	// The confirmed order is returned from API.
	resp, err := router.subsClient.QueryOrder(order)
	if err != nil {
		return render.NewInternalError(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
