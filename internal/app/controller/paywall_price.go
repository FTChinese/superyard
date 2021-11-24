package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/labstack/echo/v4"
)

func (router PaywallRouter) ListPriceOfProduct(c echo.Context) error {
	productID := c.QueryParam("product_id")
	if productID == "" {
		return render.NewBadRequest("Missing query parameter product_id")
	}

	live := getParamLive(c)

	resp, err := router.apiClients.
		Select(live).
		ListPriceOfProduct(productID)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) CreatePrice(c echo.Context) error {

	defer c.Request().Body.Close()

	live := getParamLive(c)

	resp, err := router.apiClients.
		Select(live).
		CreatePrice(c.Request().Body)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) UpdatePrice(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Param("priceId")

	live := getParamLive(c)

	resp, err := router.apiClients.
		Select(live).
		UpdatePrice(id, c.Request().Body)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) ActivatePrice(c echo.Context) error {
	id := c.Param("priceId")

	live := getParamLive(c)

	resp, err := router.apiClients.
		Select(live).
		ActivatePrice(id)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) RefreshPriceDiscounts(c echo.Context) error {
	id := c.Param("priceId")

	live := getParamLive(c)

	resp, err := router.apiClients.
		Select(live).
		RefreshPriceDiscounts(id)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) CreateDiscount(c echo.Context) error {

	live := getParamLive(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClients.
		Select(live).
		CreateDiscount(c.Request().Body)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) RemoveDiscount(c echo.Context) error {
	id := c.Param("id")

	live := getParamLive(c)

	resp, err := router.apiClients.
		Select(live).
		RemoveDiscount(id)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
