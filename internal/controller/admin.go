package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/repository/admin"
	"github.com/FTChinese/superyard/internal/repository/user"
	"github.com/FTChinese/superyard/pkg/letter"
	"github.com/FTChinese/superyard/pkg/staff"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

// AdminRouter manages staff.
type AdminRouter struct {
	postman   postoffice.PostOffice
	adminRepo admin.Env
	userRepo  user.Env
	logger    *zap.Logger
}

// NewAdminRouter creates a new instance of StaffController
func NewAdminRouter(db *sqlx.DB, p postoffice.PostOffice) AdminRouter {
	l, _ := zap.NewProduction()
	return AdminRouter{
		adminRepo: admin.NewEnv(db),
		userRepo:  user.NewEnv(db),
		postman:   p,
		logger:    l,
	}
}

// Creates creates a new account.
// Input:
// {
//	userName: string,
//  password: string
//	email: string,
//	displayName?: string,
//	department?: string,
//	groupMembers: number
// }
// Requires admin privilege.
func (router AdminRouter) CreateStaff(c echo.Context) error {
	var input staff.InputData

	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}
	if ve := input.ValidateSignUp(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	su := staff.NewSignUp(input)

	if err := router.adminRepo.CreateStaff(su); err != nil {
		return render.NewDBError(err)
	}

	go func() {
		parcel, err := letter.SignUpParcel(su, input.SourceURL)
		if err != nil {
		}
		if err := router.postman.Deliver(parcel); err != nil {
		}
	}()

	return c.JSON(http.StatusOK, su.Account)
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

	for i, p := range accounts.Data {
		if p.ID.IsZero() {
			accounts.Data[i].ID = null.StringFrom(staff.GenStaffID())
		}
	}

	go func() {
		for _, account := range accounts.Data {
			_ = router.userRepo.AddID(account)
		}
	}()

	return c.JSON(http.StatusOK, accounts)
}

// Profile shows a adminRepo's profile.
//
//	 GET /adminRepo/:id
func (router AdminRouter) StaffProfile(c echo.Context) error {

	id := c.Param("id")

	p, err := router.adminRepo.StaffProfile(id)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, p)
}

// Update updates a user's account
// Input:
// {
//	userName: string,
//	email: string,
//	displayName?: string,
//	department?: string;
//	groupMembers: number
// }
func (router AdminRouter) UpdateStaff(c echo.Context) error {

	id := c.Param("id")

	var input staff.InputData
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateAccount(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	// First retrieve current profile.
	account, err := router.adminRepo.AccountByID(id)
	if err != nil {
		return render.NewDBError(err)
	}

	account = account.Update(input)

	if err := router.adminRepo.UpdateAccount(account); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router AdminRouter) DeleteStaff(c echo.Context) error {
	id := c.Param("id")

	var vip = struct {
		Revoke bool `json:"revokeVip"`
	}{}
	if err := c.Bind(&vip); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if err := router.adminRepo.Deactivate(id); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router AdminRouter) Reinstate(c echo.Context) error {
	id := c.Param("id")

	if err := router.adminRepo.Activate(id); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Search finds an employee.
// Query parameter: q=<user name>
func (router AdminRouter) Search(c echo.Context) error {
	q := c.QueryParam("q")
	if q == "" {
		return render.NewBadRequest("Missing query parameter q")
	}

	account, err := router.adminRepo.AccountByName(q)
	if err != nil {
		return render.NewDBError(err)
	}

	if account.ID.IsZero() {
		account.ID = null.StringFrom(staff.GenStaffID())
		go func() {
			_ = router.userRepo.AddID(account)
		}()
	}

	return c.JSON(http.StatusOK, account)
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
