package subsapi

import (
	"github.com/FTChinese/superyard/pkg/config"
)

const (
	basePathPaywall    = "/paywall"
	pathRefreshPaywall = basePathPaywall + "/__refresh"
	pathProductPrices  = basePathPaywall + "/prices"
	pathPriceDiscounts = basePathPaywall + "/discounts"

	basePathStripe   = "/stripe"
	pathStripePrices = basePathStripe + "/prices?refresh=true"
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
	key            string
	baseURL        string
	sandboxBaseURL string // Deprecated
	v3BaseUrl      string // Deprecated
}

func newClient(key, baseURL string) Client {
	return Client{
		key:     key,
		baseURL: baseURL,
	}
}

// APIClients contains clients to hit various versions of API.
type APIClients struct {
	Sandbox Client
	Live    Client
	V3      Client
}

// NewAPIClients creates an APIClients.
func NewAPIClients(prod bool) APIClients {
	key := config.MustSubsAPIKey().Pick(prod)

	return APIClients{
		Sandbox: newClient(key, config.MustSubsAPISandboxBaseURL().Pick(prod)),
		Live:    newClient(key, config.MustSubsAPIV4BaseURL().Pick(prod)),
		V3: newClient(
			config.MustSubsAPIKey().Pick(true),
			config.MustSubsAPIV3BaseURL().Pick(true)),
	}
}

func (c APIClients) Select(live bool) Client {
	if live {
		return c.Live
	}

	return c.Sandbox
}
