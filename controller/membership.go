package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateMember creates a membership for an account.
//
// Input: reader.MemberInput
// expireDate: string;
// payMethod: string;
// ftcPlanId: string;
// ftcId and unionId cannot be both empty.
func (router ReaderRouter) CreateMember(c echo.Context) error {
	var input reader.MemberInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// Find the account for which the membership is created.
	a, err := router.readerRepo.FtcAccount(input.CompoundID)
	if err != nil {
		return render.NewDBError(err)
	}

	// Find the plan using the ftcPlanId field.
	plan, err := router.productsRepo.LoadPlan(input.FtcPlanID.String)
	if err != nil {
		return render.NewDBError(err)
	}
	m := input.NewMembership(a, plan)

	if err := router.readerRepo.CreateMember(m); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// LoadMember retrieves membership by either ftc uuid of wechat union id.
func (router ReaderRouter) LoadMember(c echo.Context) error {
	id := c.Param("id")

	m, err := router.readerRepo.RetrieveMember(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, m)
}

// UpdateMember modifies an existing membership.
// Input: reader.MemberInput
// expireDate: string;
// payMethod: string;
// ftcPlanId: string;
func (router ReaderRouter) UpdateMember(c echo.Context) error {
	claims := getPassportClaims(c)
	id := c.Param("id")

	var input reader.MemberInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	input.CompoundID = id

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	plan, err := router.productsRepo.LoadPlan(input.FtcPlanID.String)
	if err != nil {
		return render.NewDBError(err)
	}

	result, err := router.readerRepo.UpdateMember(input, plan)
	if err != nil {
		return render.NewDBError(err)
	}

	go func() {
		_ = router.readerRepo.SnapshotMember(
			result.Snapshot.
				WithCreator(claims.Username),
		)
	}()

	return c.JSON(http.StatusOK, result.Membership)
}
