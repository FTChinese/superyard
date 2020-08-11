package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"gitlab.com/ftchinese/superyard/pkg/paywall"
	"net/http"
)

// CreateBanner creates a single unique banner.
// Input
// heading: string;
// coverUrl?: string;
// subHeading?: string;
// content?; string;
func (router ProductRouter) CreateBanner(c echo.Context) error {
	claims := getPassportClaims(c)
	var input paywall.BannerInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	banner := paywall.NewBanner(input, claims.Username)

	if err := router.repo.CreateBanner(banner); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, banner)
}

func (router ProductRouter) LoadBanner(c echo.Context) error {

	banner, err := router.repo.LoadBanner()
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, banner)
}

func (router ProductRouter) UpdateBanner(c echo.Context) error {
	var input paywall.BannerInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	banner, err := router.repo.LoadBanner()
	if err != nil {
		return render.NewDBError(err)
	}

	banner = banner.Update(input)

	if err := router.repo.UpdateBanner(banner); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, banner)
}

// CreatePromo creates a new promo and apply it to current banner.
func (router ProductRouter) CreatePromo(c echo.Context) error {
	claims := getPassportClaims(c)
	var input paywall.PromoInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	promo := paywall.NewPromo(input, claims.Username)

	if err := router.repo.CreatePromo(promo); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, promo)
}

func (router ProductRouter) LoadPromo(c echo.Context) error {
	id := c.Param("id")

	promo, err := router.repo.LoadProduct(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, promo)
}

// ListPaywallProducts retrieves all active products presented on paywall,
// together with each product's plans, and each plan's optional
// discount.
func (router ProductRouter) ListPaywallProducts(c echo.Context) error {

	prods, err := router.repo.LoadPaywallProducts()
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, prods)
}
