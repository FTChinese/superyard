package controller

import (
	"github.com/FTChinese/go-rest/render"
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
