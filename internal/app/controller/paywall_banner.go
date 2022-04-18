package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
)

// CreateBanner creates a single unique banner.
// Request: paywall.BannerInput
// { heading: string;
//   coverUrl?: string;
//   subHeading?: string;
//   content?; string;
// }
// Response: paywall.Banner
func (routes PaywallRoutes) CreateBanner(c echo.Context) error {
	claims := getPassportClaims(c)

	live := xhttp.GetQueryLive(c)

	defer c.Request().Body.Close()

	resp, err := routes.apiClients.
		Select(live).
		CreatePaywallBanner(c.Request().Body, claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// CreatePromoBanner creates a new promo and apply it to current banner.
// Request: paywall.PromoInput
// Response: paywall.Promo
func (routes PaywallRoutes) CreatePromoBanner(c echo.Context) error {
	claims := getPassportClaims(c)

	live := xhttp.GetQueryLive(c)

	defer c.Request().Body.Close()

	resp, err := routes.apiClients.
		Select(live).
		CreatePaywallPromoBanner(
			c.Request().Body,
			claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) DropPromoBanner(c echo.Context) error {
	claims := getPassportClaims(c)
	live := xhttp.GetQueryLive(c)

	resp, err := routes.apiClients.
		Select(live).
		DropPaywallPromo(claims.Username)

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
