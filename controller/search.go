package controller

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/builder"
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/search"
	"gitlab.com/ftchinese/superyard/repository/staff"
	"net/http"
	"strings"
)

type SearchRouter struct {
	env      search.Env
	staffEnv staff.Env
}

func NewSearchRouter(db *sqlx.DB) SearchRouter {
	return SearchRouter{
		env:      search.Env{DB: db},
		staffEnv: staff.Env{DB: db},
	}
}

// SearchStaff finds a staff by email or user name
//
//	GET /search/staff?[name|email]=<value>
func (router SearchRouter) Staff(c echo.Context) error {

	p := builder.NewSearchParam(
		c.QueryParams(),
		[]string{"name", "email"},
	)

	if err := p.Validate(); err != nil {
		return util.NewBadRequest(err.Error())
	}

	col, err := employee.ParseColumn(p.Key)
	if err != nil {
		return util.NewBadRequest(err.Error())
	}

	account, err := router.env.Staff(col, p.Value)
	if err != nil {
		return util.NewDBFailure(err)
	}

	if account.ID.IsZero() {
		account.GenerateID()
		go func() {
			if err := router.staffEnv.AddID(account); err != nil {
				logger.WithField("trace", "Env.SearchStaff").Error(err)
			}
		}()
	}

	return c.JSON(http.StatusOK, account)
}

// SearchFtcUser tries to find a user by userName or email
//
//	GET /search/reader?email=<name@example.org>
func (router SearchRouter) SearchFtcUser(c echo.Context) error {
	p := builder.NewSearchParam(
		c.QueryParams(),
		[]string{"email"},
	)

	if err := p.Validate(); err != nil {
		return util.NewBadRequest(err.Error())
	}

	ftcInfo, err := router.env.SearchFtcUser(p.Value)

	// 404 Not Found
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, ftcInfo)
}

// FindWxUser tries to find a wechat user by nickname
//
// GET /search/reader/wx?q=<nickname>&page=<number>&per_page=<number>
func (router SearchRouter) SearchWxUser(c echo.Context) error {
	nickname := strings.TrimSpace(c.QueryParam("q"))

	if nickname == "" {
		return util.NewBadRequest("missing query parameter q")
	}

	var pagination builder.Pagination
	if err := c.Bind(pagination); err != nil {
		return util.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	wxUsers, err := router.env.SearchWxUser(nickname, pagination)

	// 404 Not Found
	if err != nil {
		return util.NewDBFailure(err)
	}

	return c.JSON(http.StatusOK, wxUsers)
}
