package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/labstack/echo/v4"
)

func (router B2BRouter) LoadTeam(c echo.Context) error {
	id := c.Param("id")

	resp, err := router.apiClient.LoadTeam(id)
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
