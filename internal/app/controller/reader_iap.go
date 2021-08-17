package controller

import (
	"github.com/FTChinese/go-rest/render"
	apple2 "github.com/FTChinese/superyard/internal/pkg/apple"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (router ReaderRouter) IAPMember(c echo.Context) error {
	origTxID := c.Param("id")

	m, err := router.readerRepo.IAPMember(origTxID)
	if err != nil {
		return render.NewDBError(err)
	}

	if m.IsZero() {
		return render.NewNotFound("Not found")
	}

	return c.JSON(http.StatusOK, m)
}

// LinkIAP links an existing IAP to an ftc account and creates the membership derived.
//
// POST /iap/:id/link
// { ftcId: string }
func (router ReaderRouter) LinkIAP(c echo.Context) error {
	origTxID := c.Param("id")
	var input apple2.LinkInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	input.OriginalTxID = origTxID

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	resp, errs := router.subsClient.LinkIAP(input)
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
	var input apple2.LinkInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	input.OriginalTxID = origTxID

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	resp, errs := router.subsClient.UnlinkIAP(input)
	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	// 204 no content.
	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) ListIAPSubs(c echo.Context) error {
	userID := getUserID(c)

	resp, errs := router.subsClient.ListIAPSubs(userID, c.QueryString())

	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) LoadIAPSubs(c echo.Context) error {
	id := c.Param("id")

	resp, errs := router.subsClient.LoadIAPSubs(id)

	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) RefreshIAPSubs(c echo.Context) error {
	id := c.Param("id")

	resp, errs := router.subsClient.RefreshIAPSubs(id)

	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
