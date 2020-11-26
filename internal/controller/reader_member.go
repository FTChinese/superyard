package controller

import (
	"errors"
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/letter"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/labstack/echo/v4"
	"net/http"
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

	var input subs.FtcSubsCreationInput
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
	input.PlanID = plan.ID

	account, err := router.readerRepo.CreateFtcMember(input)
	if err != nil {
		var ve *render.ValidationError
		if errors.As(err, &ve) {
			return render.NewUnprocessable(ve)
		}

		return render.NewDBError(err)
	}

	// Send membership created email.
	if account.FtcID.Valid {
		go func() {
			parcel, err := letter.MemberUpsertParcel(account)
			if err != nil {
				return
			}

			_ = router.postman.Deliver(parcel)
		}()
	}

	return c.JSON(http.StatusOK, account.Membership)
}

// LoadMember retrieves membership by either ftc uuid of wechat union id.
func (router ReaderRouter) LoadMember(c echo.Context) error {
	id := c.Param("id")

	m, err := router.readerRepo.MemberByCompoundID(id)
	if err != nil {
		return render.NewDBError(err)
	}

	if m.IsZero() {
		return render.NewNotFound("Not found")
	}

	return c.JSON(http.StatusOK, m)
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
	claims := getPassportClaims(c)
	compoundID := c.Param("id")

	var input subs.FtcSubsUpdateInput
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

	input.PlanID = plan.ID

	result, err := router.readerRepo.UpdateFtcMember(compoundID, input)
	if err != nil {
		var ve *render.ValidationError
		if errors.As(err, &ve) {
			return render.NewUnprocessable(ve)
		}

		return render.NewDBError(err)
	}

	// This is an update. Snapshot must exists.
	go func() {
		_ = router.readerRepo.SaveMemberSnapshot(
			result.Snapshot.
				WithCreator(claims.Username),
		)
	}()

	// Send membership updated email, only if the if an ftc account.
	if result.Membership.FtcID.Valid {
		go func() {
			a, err := router.readerRepo.JoinedAccountByFtcOrWx(result.Membership.IDs)
			if err != nil {
				return
			}

			parcel, err := letter.MemberUpsertParcel(reader.Account{
				JoinedAccount: a,
				Membership:    result.Membership,
			})

			if err != nil {
				return
			}
			_ = router.postman.Deliver(parcel)
		}()
	}

	return c.JSON(http.StatusOK, result.Membership)
}

// DeleteFtcMember drops membership from a user by either ftc id or union id.
func (router ReaderRouter) DeleteFtcMember(c echo.Context) error {
	claims := getPassportClaims(c)

	compoundID := c.Param("id")

	snapshot, err := router.readerRepo.DeleteFtcMember(compoundID)
	if err != nil {
		return render.NewDBError(err)
	}

	go func() {
		_ = router.readerRepo.SaveMemberSnapshot(snapshot.WithCreator(claims.Username))
	}()

	return c.NoContent(http.StatusNoContent)
}

// ListSnapshots list a user's membership revision history.
func (router ReaderRouter) ListSnapshots(c echo.Context) error {
	var page gorest.Pagination
	if err := c.Bind(&page); err != nil {
		return render.NewBadRequest(err.Error())
	}

	var ids reader.IDs
	if err := c.Bind(&ids); err != nil {
		return render.NewBadRequest(err.Error())
	}

	list, err := router.readerRepo.ListMemberSnapshots(ids, page)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, list)
}
