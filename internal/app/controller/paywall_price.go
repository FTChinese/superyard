package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
)

func (router PaywallRouter) ListPriceOfProduct(c echo.Context) error {
	productID := c.QueryParam("product_id")
	if productID == "" {
		return render.NewBadRequest("Missing query parameter product_id")
	}

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.
		Select(live).
		ListPriceOfProduct(productID, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) CreatePrice(c echo.Context) error {

	defer c.Request().Body.Close()

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.
		Select(live).
		CreatePrice(c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) UpdatePrice(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.
		Select(live).
		UpdatePrice(id, c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) ActivatePrice(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.
		Select(live).
		ActivatePrice(id, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	if live {
		go func() {
			_, err := router.apiClients.
				V5.
				ActivatePrice(id, claims.Username)

			if err != nil {
				sugar.Error(err)
			}
		}()
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) ArchivePrice(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.
		Select(live).ArchivePrice(id, claims.Username)

	if err != nil {
		_ = render.NewBadRequest(err.Error())
	}

	if live {
		go func() {
			_, err := router.apiClients.V5.ArchivePrice(
				id, claims.Username)
			if err != nil {
				sugar.Error(err)
			}
		}()
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) RefreshPriceDiscounts(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	id := c.Param("priceId")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.
		Select(live).
		RefreshPriceDiscounts(id, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	if live {
		go func() {
			_, err := router.apiClients.
				V5.
				RefreshPriceDiscounts(id, claims.Username)
			if err != nil {
				sugar.Error(err)
			}
		}()
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) CreateDiscount(c echo.Context) error {

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := router.apiClients.
		Select(live).
		CreateDiscount(c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) RemoveDiscount(c echo.Context) error {
	id := c.Param("id")

	live := xhttp.GetQueryLive(c)
	claims := getPassportClaims(c)

	resp, err := router.apiClients.
		Select(live).
		RemoveDiscount(id, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
