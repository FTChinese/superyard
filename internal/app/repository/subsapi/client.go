package subsapi

import (
	"github.com/FTChinese/superyard/pkg/config"
)

const (
	pathPaywall        = "/paywall"
	pathRefreshPaywall = pathPaywall + "/__refresh"
	pathPaywallPrices  = pathPaywall + "/active/prices"
	pathProductPrices  = pathPaywall + "/prices"
	pathPriceDiscounts = pathPaywall + "/discounts"
)

func pathPricesOfProduct(id string) string {
	return pathProductPrices + "?product_id=" + id
}

func pathPriceOf(id string) string {
	return pathProductPrices + "/" + id
}

func pathDiscountOf(id string) string {
	return pathPriceDiscounts + "/" + id
}

type Client struct {
	key     string
	baseURL string
}

func NewClient(prod bool) Client {

	return Client{
		key:     config.MustSubsAPIKey().Pick(prod),
		baseURL: config.MustSubsAPISandboxBaseURL().Pick(prod),
	}
}
