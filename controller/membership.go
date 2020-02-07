package controller

import (
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/reader"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/customer"
	"net/http"
)

type MemberRouter struct {
	env customer.Env
}

func NewMemberRouter(db *sqlx.DB) MemberRouter {
	return MemberRouter{
		env: customer.Env{DB: db},
	}
}

func (router MemberRouter) ListMembers(c echo.Context) error {
	return c.String(http.StatusOK, "Not implemented")
}

func (router MemberRouter) CreateMember(c echo.Context) error {

	var m reader.Membership
	if err := c.Bind(&m); err != nil {
		return util.NewBadRequest(err.Error())
	}

	m.GenerateID()

	if ie := m.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}

	if err := router.env.CreateMember(m); err != nil {
		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router MemberRouter) LoadMember(c echo.Context) error {

	id := c.Param("id")

	m, err := router.env.RetrieveMember(id)
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, m)
}

func (router MemberRouter) UpdateMember(c echo.Context) error {

	id := c.Param("id")

	var m reader.Membership
	if err := c.Bind(&m); err != nil {
		return util.NewBadRequest(err.Error())
	}
	m.ID = null.StringFrom(id)

	if ie := m.Validate(); ie != nil {
		return util.NewUnprocessable(ie)
	}

	if err := router.env.UpdateMember(m); err != nil {
		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router MemberRouter) DeleteMember(c echo.Context) error {

	id := c.Param("id")

	if err := router.env.DeleteMember(id); err != nil {
		return util.NewDBFailure(err)
	}

	return c.NoContent(http.StatusNoContent)
}
