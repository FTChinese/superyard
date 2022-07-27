package subsapi

import (
	"github.com/FTChinese/superyard/pkg/config"
	"github.com/FTChinese/superyard/pkg/fetch"
)

const (
	rootPathEmailAuth = "/auth/email"
	pathEmailSignUp   = rootPathEmailAuth + "/signup"

	rootPathAccount = "/account"
	pathProfile     = rootPathAccount + "/profile"
	pathAddress     = rootPathAccount + "/address"
	pathWxAccount   = rootPathAccount + "/wx"

	rootPathOrders = "/orders"

	rootPathPaywall    = "/paywall"
	pathPaywallBanner  = rootPathPaywall + "/banner"
	pathPaywallPromo   = rootPathPaywall + "/banner/promo"
	pathProducts       = rootPathPaywall + "/products"
	pathPrices         = rootPathPaywall + "/prices"
	pathPriceDiscounts = rootPathPaywall + "/discounts"
	pathRefreshPaywall = rootPathPaywall + "/__refresh"

	rootPathStripe    = "/stripe"
	pathStripePrices  = rootPathStripe + "/prices"
	pathStripeCoupons = rootPathStripe + "/coupons"

	rootPathApps        = "/apps"
	pathAndroidReleases = rootPathApps + "/android/releases"

	pathLegal = "/legal"

	rootPathCMS          = "/cms"
	pathCmsMembership    = rootPathCMS + "/memberships"
	pathCmsSnapshots     = rootPathCMS + "/snapshots"
	pathCmsAddOn         = rootPathCMS + "/addons"
	pathCmsStripeCoupons = rootPathCMS + "/stripe/coupons"
	pathCmsLegal         = rootPathCMS + "/legal"
	pathCmsAndroid       = rootPathCMS + "/android"
)

const (
	queryKeyProductID = "product_id"
	queryKeyRefresh   = "refresh"
)

func pathIntroOfProduct(base, id string) string {
	return fetch.
		NewURLBuilder(base).
		AddPath(pathProducts).
		AddPath(id).
		AddPath("intro").
		String()
}

func pathPriceOf(base, id string) string {
	return fetch.NewURLBuilder(base).
		AddPath(pathPrices).
		AddPath(id).
		String()
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
	// Used to refresh previous version of paywall data to keep backward compatible,
	V5 Client
	V4 Client
	V3 Client
}

// NewAPIClients creates an APIClients.
// When prod is false, both sandbox and live goes to localhost.
// Since localhost is always run with livemode set to false,
// you always get back sandbox data for development environment.
func NewAPIClients(prod bool) APIClients {
	keySelector := config.MustSubsAPIKey()
	key := keySelector.Pick(prod)
	prodKey := keySelector.Pick(true)

	return APIClients{
		Sandbox: newClient(key, config.MustSubsAPISandboxBaseURL().Pick(prod)),
		Live:    newClient(key, config.MustSubsAPIv6BaseURL().Pick(prod)),
		V5:      newClient(prodKey, config.MustSubsAPIv5BaseURL().Pick(true)),
		V4:      newClient(prodKey, config.MustSubsAPIV4BaseURL().Pick(true)),
		V3:      newClient(prodKey, config.MustSubsAPIV3BaseURL().Pick(true)),
	}
}

func (c APIClients) Select(live bool) Client {
	if live {
		return c.Live
	}

	return c.Sandbox
}
