package subsapi

import (
	"github.com/FTChinese/superyard/pkg/config"
)

const (
	basePathMember  = "/membership"
	basePathPaywall = "/paywall"
	basePathStripe  = "/stripe"

	pathMemberSnapshot = basePathMember + "/snapshots"
	pathMemberAddOn    = basePathMember + "/addons"

	pathPaywallBanner  = basePathPaywall + "/banner"
	pathPaywallPromo   = basePathPaywall + "/banner/promo"
	pathProducts       = basePathPaywall + "/products"
	pathPrices         = basePathPaywall + "/prices"
	pathPriceDiscounts = basePathPaywall + "/discounts"
	pathRefreshPaywall = basePathPaywall + "/__refresh"

	pathStripePrices = basePathStripe + "/prices?refresh=true"
)

func pathProductOf(id string) string {
	return pathProducts + "/" + id
}

func pathActivateProductOf(id string) string {
	return pathProducts + "/" + id + "/activate"
}

func pathPricesOfProduct(id string) string {
	return pathPrices + "?product_id=" + id
}

func pathPriceOf(id string) string {
	return pathPrices + "/" + id
}

func pathDiscountOf(id string) string {
	return pathPriceDiscounts + "/" + id
}

type Client struct {
	key     string
	baseURL string
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
// When prod is false, both sandbox and live goes to localhost.
// Since localhost is always run with livemode set to false,
// you always get back sandbox data for development environment.
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
