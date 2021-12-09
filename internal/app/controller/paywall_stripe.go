package controller

import (
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (router PaywallRouter) ListStripePrices(c echo.Context) error {
	live := xhttp.GetQueryLive(c)

	list, err := router.stripeClients.Select(live).ListPrices()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, list)
}

func (router PaywallRouter) LoadStripePrice(c echo.Context) error {
	id := c.Param("id")
	live := xhttp.GetQueryLive(c)
	p, err := router.stripeClients.
		Select(live).
		RetrievePrice(id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, p)
}
