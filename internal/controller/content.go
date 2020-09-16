package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/repository/apn"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ContentRouter struct {
	model apn.ArticleEnv
}

func NewContentRouter(db *sqlx.DB) ContentRouter {
	return ContentRouter{
		model: apn.ArticleEnv{DB: db},
	}
}

func (router ContentRouter) LatestStoryList(c echo.Context) error {
	teasers, err := router.model.LatestStoryList()

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, teasers)
}

func (router ContentRouter) StoryTeaser(c echo.Context) error {
	id := c.Param("id")

	teaser, err := router.model.FindStory(id)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, teaser)
}

func (router ContentRouter) VideoTeaser(c echo.Context) error {
	id := c.Param("id")

	teaser, err := router.model.FindVideo(id)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, teaser)
}

func (router ContentRouter) GalleryTeaser(c echo.Context) error {
	id := c.Param("id")

	teaser, err := router.model.FindGallery(id)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, teaser)
}

func (router ContentRouter) InteractiveTeaser(c echo.Context) error {
	id := c.Param("id")

	teaser, err := router.model.FindInteractive(id)

	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, teaser)
}