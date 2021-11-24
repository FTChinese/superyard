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
// Input: subs.FtcSubsCreationInput
// ftcId?: string;
// unionId?: string
// tier: string;
// cycle: string;
// expireDate: string;
// payMethod: string;
func (router ReaderRouter) CreateFtcMember(c echo.Context) error {

	resp, err := router.subsClient.CreateMembership(c.Request().Body)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// LoadMember retrieves membership by either ftc uuid of wechat union id.
func (router ReaderRouter) LoadMember(c echo.Context) error {

	resp, err := router.subsClient.LoadMembership()

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// UpdateFtcMember update or create a membership purchased via ali or wx.
//
// POST /memberships/:id
//
// Input: subs.FtcSubsUpdateInput
// tier: string;
// cycle: string;
// expireDate: string;
// payMethod: string;
func (router ReaderRouter) UpdateFtcMember(c echo.Context) error {

	resp, err := router.subsClient.UpdateMembership(c.Request().Body)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// DeleteFtcMember drops membership from a user by either ftc id or union id.
func (router ReaderRouter) DeleteFtcMember(c echo.Context) error {

	resp, err := router.
		subsClient.
		DeleteMembership(c.Request().Body)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// ListSnapshots list a user's membership revision history.
func (router ReaderRouter) ListSnapshots(c echo.Context) error {

	resp, err := router.subsClient.ListSnapshot()

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
