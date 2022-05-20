package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type LegalRoutes struct {
	apiClient subsapi.Client
	logger    *zap.Logger
}

func NewLegalRoutes(client subsapi.Client, logger *zap.Logger) LegalRoutes {
	return LegalRoutes{
		apiClient: client,
		logger:    logger,
	}
}

func (routes LegalRoutes) List(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	claims := getPassportClaims(c)

	rawQuery := c.QueryString()
	resp, err := routes.apiClient.ListLegalDocs(rawQuery, claims.Username)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}

func (routes LegalRoutes) Load(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	id := c.Param("id")

	resp, err := routes.apiClient.LoadLegalDoc(id)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}

func (routes LegalRoutes) Create(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	claims := getPassportClaims(c)

	defer c.Request().Body.Close()

	resp, err := routes.
		apiClient.
		CreateLegalDoc(c.Request().Body, claims.Username)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}

func (routes LegalRoutes) Update(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	claims := getPassportClaims(c)
	id := c.Param("id")

	defer c.Request().Body.Close()

	resp, err := routes.
		apiClient.
		UpdateLegalDoc(id, c.Request().Body, claims.Username)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}

func (routes LegalRoutes) Publish(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	claims := getPassportClaims(c)
	id := c.Param("id")

	defer c.Request().Body.Close()

	resp, err := routes.
		apiClient.
		PublishLegalDoc(id, c.Request().Body, claims.Username)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(
		resp.StatusCode,
		fetch.ContentJSON,
		resp.Body)
}
