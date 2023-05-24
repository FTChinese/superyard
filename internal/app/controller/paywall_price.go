package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
)

func (routes PaywallRoutes) ListPriceOfProduct(c echo.Context) error {
	productID := c.QueryParam("product_id")
	if productID == "" {
		return render.NewBadRequest("Missing query parameter product_id")
	}

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).
		ListPriceOfProduct(productID, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) CreatePrice(c echo.Context) error {

	defer c.Request().Body.Close()

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).
		CreatePrice(c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) LoadPrice(c echo.Context) error {
	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := routes.apiClients.Select(live).LoadFtcPrice(id, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) UpdatePrice(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).
		UpdatePrice(id, c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) ActivatePrice(c echo.Context) error {
	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).
		ActivatePrice(id, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) ArchivePrice(c echo.Context) error {
	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).ArchivePrice(id, claims.Username)

	if err != nil {
		_ = render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) RefreshPriceDiscounts(c echo.Context) error {
	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).
		RefreshPriceDiscounts(id, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) CreateDiscount(c echo.Context) error {

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := routes.apiClients.
		Select(live).
		CreateDiscount(c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) RemoveDiscount(c echo.Context) error {
	id := c.Param("id")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := routes.apiClients.
		Select(live).
		RemoveDiscount(id, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
