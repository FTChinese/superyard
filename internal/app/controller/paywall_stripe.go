package controller

import (
	"net/http"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
)

func (routes PaywallRoutes) ListStripePrices(c echo.Context) error {
	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)
	query := c.QueryParams()

	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}
	page.Normalize()

	resp, err := routes.apiClients.
		Select(live).
		ListStripePrices(query, claims.Username)

	if err != nil {
		return err
	}

	return c.Stream(http.StatusOK, fetch.ContentJSON, resp.Body)
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

func (routes PaywallRoutes) UpdateStripePriceMeta(c echo.Context) error {

	claims := getPassportClaims(c)
	priceID := c.Param("id")
	live := xhttp.GetQueryLive(c)

	defer c.Request().Body.Close()

	resp, err := routes.apiClients.
		Select(live).
		UpdateStripePriceMeta(
			priceID,
			c.Request().Body,
			claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// ActivateStripePrice put a stripe on paywall
// by insert it into subs_product.product_active_price.
func (routes PaywallRoutes) ActivateStripePrice(c echo.Context) error {

	claims := getPassportClaims(c)
	priceID := c.Param("id")
	live := xhttp.GetQueryLive(c)

	resp, err := routes.apiClients.
		Select(live).
		ActivateStripePrice(
			priceID,
			claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) DeactivateStripePrice(c echo.Context) error {

	claims := getPassportClaims(c)
	priceID := c.Param("id")
	live := xhttp.GetQueryLive(c)

	resp, err := routes.apiClients.
		Select(live).
		DeactivateStripePrice(
			priceID,
			claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
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

func (routes PaywallRoutes) ActivateCoupon(c echo.Context) error {
	id := c.Param("id")
	live := xhttp.GetQueryLive(c)

	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := routes.apiClients.
		Select(live).
		ActivateStripeCoupon(
			id,
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
