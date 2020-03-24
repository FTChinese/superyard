package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/paywall"
	"net/http"

	"gitlab.com/ftchinese/superyard/models/promo"
)

// PromoRouter handles request for subs related data.
type PromoRouter struct {
	model paywall.PromoEnv
}

// NewPromoRouter creates a new instance of SubscriptionRouter
func NewPromoRouter(db *sqlx.DB) PromoRouter {
	return PromoRouter{
		model: paywall.PromoEnv{DB: db},
	}
}

// CreateSchedule saves the schedule part of a promotion campaign.
//
//	POST /subs/schedule
//
// Input {id: number, name: string, description: null | string, startAt: string, endAt: string}
func (router PromoRouter) CreateSchedule(c echo.Context) error {
	userName := c.Request().Header.Get(userNameKey)

	var sch promo.Schedule
	if err := c.Bind(&sch); err != nil {
		return render.NewBadRequest(err.Error())
	}

	sch.Sanitize()

	if ve := sch.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	id, err := router.model.NewSchedule(sch, userName)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"id": id,
	})
}

// SetPricingPlans saves/updates a promotion's pricing plans.
//
// PATCH /subs/schedule/:id/pricing
func (router PromoRouter) SetPricingPlans(c echo.Context) error {
	id, err := ParseInt(c.Param("id"))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	var plans promo.Pricing
	if err := c.Bind(&plans); err != nil {
		return render.NewBadRequest(err.Error())
	}

	err = router.model.SavePlans(id, plans)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// SetPromoBanner saves/updates a promotion's banner content
//
// POST /subs/schedule/:id/banner
func (router PromoRouter) SetBanner(c echo.Context) error {
	id, err := ParseInt(c.Param("id"))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	var banner promo.Banner
	if err := c.Bind(&banner); err != nil {
		return render.NewBadRequest(err.Error())
	}

	banner.Sanitize()
	if ve := banner.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	err = router.model.SaveBanner(id, banner)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// ListPromos list promotion schedules by page.
//
// GET /subs/promos?page=<int>&per_page=<number>
func (router PromoRouter) ListPromos(c echo.Context) error {

	var pagination util.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	promos, err := router.model.ListPromos(pagination)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, promos)
}

// GetPromo loads a piece of promotion.
//
// GET /subs/promos/:id
func (router PromoRouter) LoadPromo(c echo.Context) error {
	id, err := ParseInt(c.Param("id"))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	p, err := router.model.LoadPromo(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, p)
}

// RemovePromo deletes a record.
//
// DELETE `/subs/promos/:id`
func (router PromoRouter) DisablePromo(c echo.Context) error {
	id, err := ParseInt(c.Param("id"))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	err = router.model.DisablePromo(id)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}
