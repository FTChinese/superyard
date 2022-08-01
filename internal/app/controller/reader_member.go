package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/labstack/echo/v4"
)

// CreateFtcMember update or create a membership purchased via ali or wx.
//
// POST /memberships
//
// - ftcId?: string;
// - unionId?: string;
// - tier: string;
// - cycle: string;
// - expireDate: string;
// - payMethod: string;
func (router ReaderRouter) CreateFtcMember(c echo.Context) error {

	claims := getPassportClaims(c)

	resp, err := router.APIClients.
		Select(true).
		CreateMembership(
			c.Request().Body,
			claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// DeleteFtcMember drops membership from a user by either ftc id or union id.
func (router ReaderRouter) DeleteFtcMember(c echo.Context) error {

	claims := getPassportClaims(c)
	id := c.Param("id")

	resp, err := router.
		APIClients.
		Select(true).
		DeleteMembership(
			id,
			c.Request().Body,
			claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// ListSnapshots list a user's membership revision history.
func (router ReaderRouter) ListSnapshots(c echo.Context) error {

	claims := getPassportClaims(c)
	query := c.QueryParams()

	resp, err := router.APIClients.
		Select(true).
		ListSnapshot(query, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
