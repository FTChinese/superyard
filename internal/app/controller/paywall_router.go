package controller

import (
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/fetch"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strconv"
)

type PaywallRouter struct {
	apiClients subsapi.APIClients
	logger     *zap.Logger
}

func NewPaywallRouter(clients subsapi.APIClients, logger *zap.Logger) PaywallRouter {
	return PaywallRouter{
		apiClients: clients,
		logger:     logger,
	}
}

// LoadPaywall gets a paywall's banner, optional promo and a list of products.
func (router PaywallRouter) LoadPaywall(c echo.Context) error {

	liveMode := getParamLive(c)

	resp, err := router.apiClients.Select(liveMode).LoadPaywall()

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
func (router PaywallRouter) RefreshFtcPaywall(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	liveMode := getParamLive(c)

	resp, err := router.apiClients.Select(liveMode).RefreshFtcPaywall()
	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	// Also bust cache of v3.
	if liveMode {
		go func() {
			sugar.Infof("Paywall cach bust v3")
			_, err := router.apiClients.V3.RefreshFtcPaywall()
			if err != nil {
				sugar.Error(err)
			}
		}()

		go func() {
			sugar.Infof("Paywall cache bust v3")
			_, err := router.apiClients.V4.RefreshFtcPaywall()
			if err != nil {
				sugar.Error(err)
			}
		}()
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}

func (router PaywallRouter) RefreshStripePrices(c echo.Context) error {
	defer router.logger.Sync()
	sugar := router.logger.Sugar()

	liveMode, _ := strconv.ParseBool(c.QueryParam("live"))

	resp, err := router.apiClients.Select(liveMode).RefreshStripePrices()
	if err != nil {
		sugar.Error(err)
		return render.NewBadRequest(err.Error())
	}

	// Also bust cache of v3.
	if liveMode {
		go func() {
			sugar.Infof("Paywall cach bust v3")
			_, err := router.apiClients.V3.RefreshStripePrices()
			if err != nil {
				sugar.Error(err)
			}
		}()
	}

	return c.Stream(resp.StatusCode, fetch.ContentJSON, resp.Body)
}
