package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/pkg/apple"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
)

// LinkIAP links an existing IAP to an ftc account and creates the membership derived.
//
// POST /iap/:id/link
// { ftcId: string }
func (router ReaderRouter) LinkIAP(c echo.Context) error {
	origTxID := c.Param("id")
	var input apple.LinkInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	input.OriginalTxID = origTxID

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	resp, errs := router.APIClients.Select(true).LinkIAP(input)
	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// UnlinkIAP severs the links between IAP and ftc account.
//
// POST /iap/:id/unlink
// { ftcId: string }
func (router ReaderRouter) UnlinkIAP(c echo.Context) error {
	origTxID := c.Param("id")
	var input apple.LinkInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	input.OriginalTxID = origTxID

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	resp, errs := router.APIClients.Select(true).UnlinkIAP(input)
	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	// 204 no content.
	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) ListIAPSubs(c echo.Context) error {
	userID := xhttp.GetFtcID(c)

	resp, errs := router.APIClients.Select(true).ListIAPSubs(userID, c.QueryString())

	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) LoadIAPSubs(c echo.Context) error {
	id := c.Param("id")

	resp, errs := router.APIClients.Select(true).LoadIAPSubs(id)

	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) RefreshIAPSubs(c echo.Context) error {
	id := c.Param("id")

	resp, errs := router.APIClients.Select(true).RefreshIAPSubs(id)

	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
