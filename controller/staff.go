package controller

import (
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/models/employee"
	"gitlab.com/ftchinese/superyard/models/util"
	"gitlab.com/ftchinese/superyard/repository/admin"
	"gitlab.com/ftchinese/superyard/repository/user"
	"net/http"
)

// StaffRouter responds to CMS login and personal settings.
type StaffRouter struct {
	postman   postoffice.Postman
	adminRepo admin.Env
	userRepo  user.Env
}

// NewStaffRouter creates a new instance of StaffController
func NewStaffRouter(db *sqlx.DB, p postoffice.Postman) StaffRouter {
	return StaffRouter{
		adminRepo: admin.Env{DB: db},
		userRepo:  user.Env{DB: db},
		postman:   p,
	}
}

// Creates creates a new account.
// Input:
// {
//	userName: string,
//	email: string,
//	displayName?: string,
//	department?: string,
//	groupMembers: number
// }
// Requires admin privilege.
func (router StaffRouter) Create(c echo.Context) error {
	var a employee.BaseAccount

	if err := c.Bind(&a); err != nil {
		return render.NewBadRequest(err.Error())
	}
	if ve := a.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	su := employee.NewSignUp(a)

	if err := router.adminRepo.Create(su); err != nil {
		return render.NewDBError(err)
	}

	go func() {
		parcel, err := su.SignUpParcel()
		if err != nil {
			logger.WithField("trace", "StaffRouter.Login").Error(err)
		}
		if err := router.postman.Deliver(parcel); err != nil {
			logger.WithField("trace", "StaffRouter.Login").Error(err)
		}
	}()

	return c.JSON(http.StatusOK, a)
}

// ListStaff shows all adminRepo.
func (router StaffRouter) List(c echo.Context) error {

	var pagination util.Pagination
	if err := c.Bind(&pagination); err != nil {
		return render.NewBadRequest(err.Error())
	}
	pagination.Normalize()

	logger.Infof("Pagination: %+v", pagination)

	accounts, err := router.adminRepo.ListStaff(pagination)
	if err != nil {
		return render.NewDBError(err)
	}

	for i, p := range accounts {
		if p.ID.IsZero() {
			accounts[i].ID = null.StringFrom(employee.GenStaffID())
		}
	}

	go func() {
		for _, account := range accounts {
			_ = router.userRepo.AddID(account)
		}
	}()

	return c.JSON(http.StatusOK, accounts)
}

// Profile shows a adminRepo's profile.
//
//	 GET /adminRepo/:id
func (router StaffRouter) Profile(c echo.Context) error {

	id := c.Param("id")

	logger.Infof("Profile for adminRepo: %s", id)

	p, err := router.adminRepo.StaffProfile(id)

	// `404 Not Found` if this user does not exist.
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, p)
}

// Update updates a user's account
func (router StaffRouter) Update(c echo.Context) error {
	log := logger.WithField("trace", "StaffRouter.Update")

	id := c.Param("id")

	// First retrieve current profile.
	account, err := router.adminRepo.AccountByID(id)
	if err != nil {
		return render.NewDBError(err)
	}

	var ba employee.BaseAccount
	if err := c.Bind(&ba); err != nil {
		log.Error(err)
		return render.NewBadRequest(err.Error())
	}

	if ve := ba.Validate(); ve != nil {
		log.Error(ve)
		return render.NewUnprocessable(ve)
	}

	account.BaseAccount = ba

	if err := router.adminRepo.UpdateAccount(account); err != nil {
		log.Error(err)
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (router StaffRouter) Delete(c echo.Context) error {
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

func (router StaffRouter) Reinstate(c echo.Context) error {
	id := c.Param("id")

	if err := router.adminRepo.Activate(id); err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// Search finds an employee.
// Query parameter: q=<user name>
func (router StaffRouter) Search(c echo.Context) error {
	q := c.QueryParam("q")
	if q == "" {
		return render.NewBadRequest("Missing query parameter q")
	}

	account, err := router.adminRepo.AccountByName(q)
	if err != nil {
		return render.NewDBError(err)
	}

	if account.ID.IsZero() {
		account.ID = null.StringFrom(employee.GenStaffID())
		go func() {
			_ = router.userRepo.AddID(account)
		}()
	}

	return c.JSON(http.StatusOK, account)
}
