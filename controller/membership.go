package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/FTChinese/superyard/repository/readers"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MemberRouter struct {
	env readers.Env
}

func NewMemberRouter(db *sqlx.DB) MemberRouter {
	return MemberRouter{
		env: readers.Env{DB: db},
	}
}

// CreateMember creates a membership for an account.
// Input:
// {
//	compoundId: string,
//	ftcId?: string,
//	unionId?: string,
//	tier: 'standard' | 'premium',
//	cycle: 'year' | 'month',
//	expireDate: string,
//	payMethod: string,
// }
// ftcId and unionId cannot be both empty.
func (router MemberRouter) CreateMember(c echo.Context) error {

	var m subs.Membership
	if err := c.Bind(&m); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := m.ValidateCreate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	if err := router.env.CreateMember(m); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router MemberRouter) LoadMember(c echo.Context) error {
	id := c.Param("id")

	m, err := router.env.RetrieveMember(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, m)
}

// UpdateMember modifies an existing membership.
// Input:
// {
//	tier: 'standard' | 'premium',
//	cycle: 'year' | 'month',
//	expireDate: string,
//	payMethod: string,
// }
func (router MemberRouter) UpdateMember(c echo.Context) error {
	claims := getPassportClaims(c)
	id := c.Param("id")

	var m subs.Membership
	if err := c.Bind(&m); err != nil {
		return render.NewBadRequest(err.Error())
	}
	m.CompoundID = id

	if ve := m.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	if err := router.env.UpdateMember(m, claims.Username); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// FindMemberForOrder tries to find membership by an order's compound id.
// An order's compound id could be either ftc id or union id.
// We cannot be sure which kind it is since there are not enough data to
// distinguish it.
func (router MemberRouter) FindMemberForOrder(c echo.Context) error {
	id := c.Param("id")

	m, err := router.env.FindMemberForOrder(id)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, m)
}
