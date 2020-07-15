package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/pkg/wiki"
	"gitlab.com/ftchinese/superyard/repository/wikis"
	"net/http"
	"strconv"
)

type WikiRouter struct {
	repo wikis.Env
}

func NewWikiRouter(db *sqlx.DB) WikiRouter {
	return WikiRouter{repo: wikis.NewEnv(db)}
}

func (router WikiRouter) CreateArticle(c echo.Context) error {
	var a wiki.Article

	if err := c.Bind(&a); err != nil {
		return render.NewBadRequest(err.Error())
	}

	id, err := router.repo.CreateArticle(a)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"id": id,
	})
}

func (router WikiRouter) UpdateArticle(c echo.Context) error {
	var a wiki.Article

	if err := c.Bind(&a); err != nil {
		return render.NewBadRequest(err.Error())
	}

	err := router.repo.UpdateArticle(a)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router WikiRouter) OneArticle(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	a, err := router.repo.LoadArticle(int64(id))
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, a)
}

func (router WikiRouter) ListArticle(c echo.Context) error {
	var p gorest.Pagination
	if err := c.Bind(&p); err != nil {
		return render.NewBadRequest(err.Error())
	}

	articles, err := router.repo.ListArticles(p)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, articles)
}
