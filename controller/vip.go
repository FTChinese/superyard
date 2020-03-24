package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/util"
	"net/http"
)

// ListVIP lists all ftc account granted vip.
//
//	GET /vip?page=<number>&per_page=<number>
func (router ReaderRouter) ListVIP(c echo.Context) error {

	var pagination util.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	vips, err := router.env.ListVIP(pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, vips)
}

// GrantVIP grants vip to an ftc account.
//
//	PUT /vip/:id
func (router ReaderRouter) GrantVIP(c echo.Context) error {
	id := c.Param("id")

	if err := router.env.GrantVIP(id); err != nil {
		return render.NewDBError(err)
	}

	// 204 No Content
	return c.NoContent(http.StatusNoContent)
}

// RevokeVIP removes a ftc account from vip.
//
//	DELETE /vip/:id
func (router ReaderRouter) RevokeVIP(c echo.Context) error {
	id := c.Param("id")

	if err := router.env.RevokeVIP(id); err != nil {
		return render.NewDBError(err)
	}

	// 204 No Content
	return c.NoContent(http.StatusNoContent)
}
