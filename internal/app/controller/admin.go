package controller

import (
	"net/http"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/admin"
	"github.com/FTChinese/superyard/internal/app/repository/auth"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/postman"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// AdminRouter manages staff.
type AdminRouter struct {
	postman   postman.Postman
	adminRepo admin.Env
	userRepo  auth.Env
	logger    *zap.Logger
}

// NewAdminRouter creates a new instance of StaffController
func NewAdminRouter(myDBs db.ReadWriteMyDBs, gormDBs db.MultiGormDBs, p postman.Postman) AdminRouter {
	l, _ := zap.NewProduction()
	return AdminRouter{
		adminRepo: admin.NewEnv(myDBs),
		userRepo:  auth.NewEnv(myDBs, gormDBs),
		postman:   p,
		logger:    l,
	}
}

// ListStaff shows all adminRepo.
func (router AdminRouter) ListStaff(c echo.Context) error {

	var pagination gorest.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	accounts, err := router.adminRepo.ListStaff(pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, accounts)
}

func (router AdminRouter) ListVIPs(c echo.Context) error {
	var p gorest.Pagination
	if err := c.Bind(&p); err != nil {
		return render.NewBadRequest(err.Error())
	}
	p.Normalize()

	vips, err := router.adminRepo.ListVIP(p)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, vips)
}

func (router AdminRouter) SetVIP(vip bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer router.logger.Sync()
		sugar := router.logger.Sugar()

		id := c.Param("id")

		a, err := router.adminRepo.FtcAccount(id)
		if err != nil {
			sugar.Error(err)
			return render.NewDBError(err)
		}

		if a.VIP == vip {
			return c.JSON(http.StatusOK, a)
		}

		a.VIP = vip

		err = router.adminRepo.UpdateVIP(a)
		if err != nil {
			sugar.Error(err)
			return render.NewDBError(err)
		}

		return c.JSON(http.StatusOK, a)
	}
}
