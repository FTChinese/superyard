package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/labstack/echo/v4"
	"net/http"
)

// CreateBanner creates a single unique banner.
// Request: paywall.BannerInput
// { heading: string;
//   coverUrl?: string;
//   subHeading?: string;
//   content?; string;
// }
// Response: paywall.Banner
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
	// Retrieve by a fixed id.
	banner, err := router.repo.LoadBanner(1)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, banner)
}

// UpdateBanner modifies the content of a banner.
// Request: paywall.BannerInput
// Response: paywall.Banner
func (router ProductRouter) UpdateBanner(c echo.Context) error {
	var input paywall.BannerInput
	if err := c.Bind(&input); err != nil {
		return render.NewBadRequest(err.Error())
	}

	if ve := input.Validate(); ve != nil {
		return render.NewUnprocessable(ve)
	}

	banner, err := router.repo.LoadBanner(1)
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
// Request: paywall.PromoInput
// Response: paywall.Promo
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

	// Create promo and apply it to banner. Since there is only one record in banner,
	// its id is fixed to 1.
	if err := router.repo.CreatePromo(1, promo); err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, promo)
}

func (router ProductRouter) LoadPromo(c echo.Context) error {
	id := c.Param("id")

	promo, err := router.repo.LoadPromo(id)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, promo)
}

// DropBannerPromo removes promo id from a banner.
// Request data: empty
// Response: 204
func (router ProductRouter) DropBannerPromo(c echo.Context) error {
	err := router.repo.DropBannerPromo(1)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// LoadPaywall gets a paywall's banner, optional promo and a list of products.
func (router ProductRouter) LoadPaywall(c echo.Context) error {
	pw, err := router.repo.LoadPaywall(1)
	if err != nil {
		return render.NewDBError(err)
	}

	return c.JSON(http.StatusOK, pw)
}

func (router ProductRouter) RefreshAPI(c echo.Context) error {
	resp, err := router.apiClient.RefreshPaywall()
	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, "application/json; charset=utf-8", resp.Body)
}
