package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	wikis2 "github.com/FTChinese/superyard/internal/app/repository/wikis"
	wiki2 "github.com/FTChinese/superyard/internal/pkg/wiki"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type WikiRouter struct {
	repo wikis2.Env
}

func NewWikiRouter(myDBs db.ReadWriteMyDBs) WikiRouter {
	return WikiRouter{repo: wikis2.NewEnv(myDBs)}
}

// Input
// {
//	title: string,
//  summary?: string,
//  keyword?: string,
//  body?: string
// }
func (router WikiRouter) CreateArticle(c echo.Context) error {
	claims := getPassportClaims(c)

	var input wiki2.ArticleInput

	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	a := wiki2.NewArticle(input, claims.Username)

	id, err := router.repo.CreateArticle(a)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"id": id,
	})
}

// UpdateArticle update an article.
// Input:
// {
//	title: string,
//  summary: string,
//  keyword?: string,
//  body: string
// }
func (router WikiRouter) UpdateArticle(c echo.Context) error {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	var input wiki2.ArticleInput

	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	article := input.Update(id)

	err = router.repo.UpdateArticle(article)
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

	p.Normalize()

	articles, err := router.repo.ListArticles(p)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, articles)
}
