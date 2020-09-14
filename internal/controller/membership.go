package controller

import (
	"errors"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/apple"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/labstack/echo/v4"
	"net/http"
)

// LoadMember retrieves membership by either ftc uuid of wechat union id.
func (router ReaderRouter) LoadMember(c echo.Context) error {
	id := c.Param("id")

	m, err := router.readerRepo.RetrieveMember(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, m)
}

// DeleteMember drops membership from a user by either ftc id or union id.
func (router ReaderRouter) DeleteMember(c echo.Context) error {
	claims := getPassportClaims(c)

	compoundID := c.Param("id")

	snapshot, err := router.readerRepo.DeleteMember(compoundID)
	if err != nil {
		return render.NewDBError(err)
	}

	go func() {
		_ = router.readerRepo.SnapshotMember(snapshot.WithCreator(claims.Username))
	}()

	return c.NoContent(http.StatusNoContent)
}

// UpsertFtcSubs update or create a membership purchased via ali or wx.
//
// PATCH /memberships/:id/ftc
//
// Input: subs.FtcSubsInput
// ftcId?: string;
// unionId?: string
// tier: string;
// cycle: string;
// expireDate: string;
// payMethod: string;
func (router ReaderRouter) UpsertFtcSubs(c echo.Context) error {
	claims := getPassportClaims(c)

	var input subs.FtcSubsInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// Get the plan the updated membership is subscribed to.
	plan, err := router.productsRepo.PaywallPlanByEdition(input.Edition)
	if err != nil {
		return render.NewDBError(err)
	}

	result, err := router.readerRepo.UpsertFtcSubs(input, plan)
	if err != nil {
		var ve *render.ValidationError
		if errors.As(err, &ve) {
			return render.NewUnprocessable(ve)
		}

		return render.NewDBError(err)
	}

	if !result.Snapshot.IsZero() {
		go func() {
			_ = router.readerRepo.SnapshotMember(
				result.Snapshot.
					WithCreator(claims.Username),
			)
		}()
	}

	return c.JSON(http.StatusOK, result.Membership)
}

// LinkIAP refreshes an existing apple subscription by original transaction id and then
// link it to an ftc account.
//
// PATCH /memberships/:id/apple
func (router ReaderRouter) LinkIAP(c echo.Context) error {
	id := c.Param("id")
	var input apple.LinkInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	input.FtcID = id

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	resp, errs := router.subsClient.LinkIAP(input)
	if errs != nil {
		return render.NewInternalError(errs[0].Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router ReaderRouter) ListIAPSubs(c echo.Context) error {
	resp, errs := router.subsClient.ListIAPSubs(c.QueryString())

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

func (router ReaderRouter) UpsertStripeSubs(c echo.Context) error {
	return render.NewInternalError("not implemented")
}
