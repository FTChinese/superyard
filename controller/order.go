package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/letter"
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

// ConfirmOrder set an order confirmation time,
// and create/renew/upgrade membership based on this order.
func (router ReaderRouter) ConfirmOrder(c echo.Context) error {
	orderID := c.Param("id")

	result, err := router.readerRepo.ConfirmOrder(orderID)

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

	// Back previous membership.
	go func() {
		if !result.Snapshot.IsZero() {
			_ = router.readerRepo.SnapshotMember(result.Snapshot)
		}
	}()

	// Send email
	go func() {
		if result.Membership.FtcID.IsZero() {
			return
		}

		account, err := router.readerRepo.FtcAccount(result.Membership.FtcID.String)
		if err != nil {
			return
		}

		parcel, err := letter.OrderConfirmedParcel(account, result)
		if err != nil {
			return
		}

		_ = router.postman.Deliver(parcel)
	}()

	return c.JSON(http.StatusOK, result.Order)
}
