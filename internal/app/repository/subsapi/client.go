package subsapi

import (
	"github.com/FTChinese/superyard/pkg/config"
	"strings"
)

const (
	rootPathMember  = "/membership"
	rootPathPaywall = "/paywall"
	rootPathStripe  = "/stripe"

	pathMemberSnapshot = rootPathMember + "/snapshots"
	pathMemberAddOn    = rootPathMember + "/addons"

	pathPaywallBanner  = rootPathPaywall + "/banner"
	pathPaywallPromo   = rootPathPaywall + "/banner/promo"
	pathProducts       = rootPathPaywall + "/products"
	pathPrices         = rootPathPaywall + "/prices"
	pathPriceDiscounts = rootPathPaywall + "/discounts"
	pathRefreshPaywall = rootPathPaywall + "/__refresh"

	pathStripePrices = rootPathStripe + "/prices?refresh=true"
)

const (
	queryKeyProductID = "product_id"
)

func pathProductOf(id string) string {
	return strings.Join([]string{pathProducts}, id)
}

func pathActivateProductOf(id string) string {
	return strings.Join([]string{pathProducts, id, "activate"}, "/")
}

func pathPriceOf(id string) string {
	return strings.Join([]string{pathPrices, id}, "/")
}

func pathActivatePriceOf(id string) string {
	return strings.Join([]string{pathPrices, id, "activate"}, "/")
}

func pathRefreshOffersOfPrice(id string) string {
	return strings.Join([]string{pathPrices, id, "discounts"}, "/")
}

func pathDiscountOf(id string) string {
	return strings.Join([]string{pathPriceDiscounts, id}, "/")
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
