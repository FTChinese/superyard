package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (routes PaywallRoutes) ListStripePrices(c echo.Context) error {
	live := xhttp.GetQueryLive(c)

	list, err := routes.apiClients.
		Select(live).
		ListStripePrices(false)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, list)
}

// LoadStripePrice gets a coupon from API.
// Query parameters:
// - live=true|false
// - refresh=true|false
func (routes PaywallRoutes) LoadStripePrice(c echo.Context) error {
	id := c.Param("id")

	var q LiveRefresh
	err := xhttp.DecodeForm(&q, c.Request())
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	resp, err := routes.apiClients.
		Select(q.Live).
		LoadStripePrice(id, q.Refresh)

	if err != nil {
		return err
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// ListStripeCoupons for a price.
// Query parameters:
// - price_id=<stripe price id>
func (routes PaywallRoutes) ListStripeCoupons(c echo.Context) error {
	priceID := c.Param("id")
	live := xhttp.GetQueryLive(c)

	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).
		ListStripePriceCoupons(priceID, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// LoadStripeCoupon gets a coupon from API.
// Query parameters:
// - live=true|false
// - refresh=true|false
func (routes PaywallRoutes) LoadStripeCoupon(c echo.Context) error {
	id := c.Param("id")

	var q LiveRefresh
	err := xhttp.DecodeForm(&q, c.Request())
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	resp, err := routes.apiClients.
		Select(q.Live).
		LoadStripeCoupon(id, q.Refresh)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) UpdateCoupon(c echo.Context) error {
	id := c.Param("id")
	live := xhttp.GetQueryLive(c)

	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := routes.apiClients.
		Select(live).
		UpdateStripeCoupon(
			id,
			c.Request().Body,
			claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) DeleteCoupon(c echo.Context) error {
	id := c.Param("id")
	live := xhttp.GetQueryLive(c)

	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).
		DeleteCoupon(
			id,
			claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
