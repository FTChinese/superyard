package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/stripeapi"
	"github.com/FTChinese/superyard/internal/app/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PaywallRoutes struct {
	apiClients    subsapi.APIClients
	stripeClients stripeapi.Clients
	logger        *zap.Logger
}

func NewPaywallRouter(clients subsapi.APIClients, logger *zap.Logger) PaywallRoutes {
	return PaywallRoutes{
		apiClients:    clients,
		stripeClients: stripeapi.NewClients(logger),
		logger:        logger,
	}
}

// LoadPaywall gets a paywall's banner, optional promo and a list of products.
func (routes PaywallRoutes) LoadPaywall(c echo.Context) error {

	liveMode := xhttp.GetQueryLive(c)

	resp, err := routes.apiClients.Select(liveMode).LoadPaywall()

	if err != nil {
		return render.NewBadRequest(err.Error())
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

// RefreshFtcPaywall bust cache of paywall, either in live mode or not.
// When busting cache, we have to clean caches for
// * Live version
// * Sandbox version
// * V3
// Plus Stripe prices, we have a total of 6 endpoints to hit.
func (routes PaywallRoutes) RefreshFtcPaywall(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	liveMode := xhttp.GetQueryLive(c)

	resp, err := routes.apiClients.Select(liveMode).RefreshFtcPaywall()
	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	// Also bust cache of v3.
	if liveMode {
		go func() {
			sugar.Infof("Paywall cach bust v3")
			_, err := routes.apiClients.V3.RefreshFtcPaywall()
			if err != nil {
				sugar.Error(err)
			}
		}()

		go func() {
			sugar.Infof("Paywall cache bust v3")
			_, err := routes.apiClients.V4.RefreshFtcPaywall()
			if err != nil {
				sugar.Error(err)
			}
		}()
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (routes PaywallRoutes) RefreshStripePrices(c echo.Context) error {
	defer routes.logger.Sync()
	sugar := routes.logger.Sugar()

	liveMode := xhttp.GetQueryLive(c)

	resp, err := routes.apiClients.
		Select(liveMode).
		ListStripePrices(true)

	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	// Also bust cache of v3.
	if liveMode {
		go func() {
			sugar.Infof("Paywall cach bust v3")
			_, err := routes.apiClients.V3.ListStripePrices(true)
			if err != nil {
				sugar.Error(err)
			}
		}()
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
