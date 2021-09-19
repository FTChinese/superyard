package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/labstack/echo/v4"
)

func (router B2BRouter) ListOrders(c echo.Context) error {
	rawQuery := c.QueryString()

	resp, err := router.apiClient.ListOrders(rawQuery)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router B2BRouter) LoadOrder(c echo.Context) error {
	id := c.Param("id")
	resp, err := router.apiClient.LoadOrder(id)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router B2BRouter) ConfirmOrder(c echo.Context) error {
	id := c.Param("id")
	resp, err := router.apiClient.ConfirmOrder(id)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
