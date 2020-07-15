package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/postoffice"
	"github.com/FTChinese/go-rest/render"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/pkg/letter"
	"gitlab.com/ftchinese/superyard/pkg/staff"
	"gitlab.com/ftchinese/superyard/repository/admin"
	"gitlab.com/ftchinese/superyard/repository/user"
	"net/http"
)

// StaffRouter responds to CMS login and personal settings.
type StaffRouter struct {
	postman   postoffice.PostOffice
	adminRepo admin.Env
	userRepo  user.Env
}

// NewStaffRouter creates a new instance of StaffController
func NewStaffRouter(db *sqlx.DB, p postoffice.PostOffice) StaffRouter {
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
//  password: string
//	email: string,
//	displayName?: string,
//	department?: string,
//	groupMembers: number
// }
// Requires admin privilege.
func (router StaffRouter) Create(c echo.Context) error {
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
			logger.WithField("trace", "StaffRouter.Login").Error(err)
		}
		if err := router.postman.Deliver(parcel); err != nil {
			logger.WithField("trace", "StaffRouter.Login").Error(err)
		}
	}()

	return c.JSON(http.StatusOK, su.Account)
}

// ListStaff shows all adminRepo.
func (router StaffRouter) List(c echo.Context) error {

	var pagination gorest.Pagination
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
			accounts[i].ID = null.StringFrom(staff.GenStaffID())
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
// Input:
// {
//	userName: string,
//	email: string,
//	displayName?: string,
//	department?: string;
//	groupMembers: number
// }
func (router StaffRouter) Update(c echo.Context) error {
	log := logger.WithField("trace", "StaffRouter.Update")

	id := c.Param("id")

	var input staff.InputData
	if err := c.Bind(&input); err != nil {
		log.Error(err)
		return render.NewBadRequest(err.Error())
	}

	if ve := input.ValidateAccount(); ve != nil {
		log.Error(ve)
		return render.NewUnprocessable(ve)
	}

	// First retrieve current profile.
	account, err := router.adminRepo.AccountByID(id)
	if err != nil {
		return render.NewDBError(err)
	}

	account = account.Update(input)

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
		account.ID = null.StringFrom(staff.GenStaffID())
		go func() {
			_ = router.userRepo.AddID(account)
		}()
	}

	return c.JSON(http.StatusOK, account)
}
