package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/app/controller"
	"github.com/FTChinese/superyard/internal/app/repository/readers"
	"github.com/FTChinese/superyard/internal/app/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/config"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/pkg/postman"
	"github.com/FTChinese/superyard/pkg/xhttp"
	"github.com/FTChinese/superyard/web"
	"github.com/flosch/pongo2/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed build/api.toml
var tomlConfig string

//go:embed build/version
var version string

//go:embed build/build_time
var build string

//go:embed client_version_next
var clientVersionNext string

var clientVersionNg string

var (
	isProduction bool
)

func init() {
	flag.BoolVar(&isProduction, "production", false, "Indicate productions environment if present")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	config.MustSetupViper([]byte(tomlConfig))
}

func newNgFooter() web.Footer {
	return web.Footer{
		Year:          time.Now().Year(),
		ClientVersion: clientVersionNg,
		ServerVersion: version,
	}
}

func NewNextFooter() web.Footer {
	return web.Footer{
		Year:          time.Now().Year(),
		ClientVersion: clientVersionNext,
		ServerVersion: version,
	}
}

func main() {
	webCfg := web.Config{
		Debug:   !isProduction,
		Version: version,
		BuiltAt: build,
	}

	logger := config.MustGetLogger(isProduction)

	gormDBs := db.MustNewMultiGormDBs(isProduction)

	ftcPm := postman.New(config.MustGetEmailConn())
	hanqiPm := postman.New(config.MustGetHanqiConn())

	appKey := config.MustGetAppKey()

	guard := controller.NewAuthGuard(appKey.GetJWTKey())

	e := echo.New()
	e.Renderer = web.MustNewRenderer(webCfg)

	if !isProduction {
		e.Static("/static", "build/public/static")
	}

	e.Pre(middleware.AddTrailingSlash())
	e.HTTPErrorHandler = errorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CSRF())
	e.Use(xhttp.DumpRequest)

	e.GET("/", func(context echo.Context) error {
		return context.Redirect(http.StatusFound, "/ng")
	})

	e.GET("/ng/*", func(c echo.Context) error {
		return c.Render(http.StatusOK, "ng.html", pongo2.Context{
			"footer": newNgFooter(),
		})
	}, xhttp.NoCache)

	e.GET("/next/*", func(c echo.Context) error {
		return c.Render(http.StatusOK, "next.html", pongo2.Context{
			"footer": NewNextFooter(),
		})
	}, xhttp.NoCache)

	apiGroup := e.Group("/api")

	apiClients := subsapi.NewAPIClients(isProduction)

	readerRouter := controller.ReaderRouter{
		Repo:       readers.New(gormDBs, logger),
		Postman:    hanqiPm,
		APIClients: apiClients,
		Logger:     logger,
		Version:    version,
	}
	productRoutes := controller.NewPaywallRouter(apiClients, logger)

	userRouter := controller.NewUserRouter(gormDBs, ftcPm, guard)

	authGroup := apiGroup.Group("/auth")
	{
		// Login
		authGroup.POST("/login/", userRouter.Login)
		// Password reset
		authGroup.POST("/password-reset/", userRouter.ResetPassword)
		authGroup.POST("/password-reset/letter/", userRouter.ForgotPassword)
		authGroup.GET("/password-reset/tokens/:token/", userRouter.VerifyResetToken)
	}

	settingsGroup := apiGroup.Group("/settings", guard.RequireLoggedIn)
	{
		// Use to renew Json Web Token
		settingsGroup.GET("/account/", userRouter.Account, guard.RequireLoggedIn)
		// Set email if empty. User can only set
		// it once.
		settingsGroup.PATCH("/account/email/", userRouter.SetEmail)
		// Allow user to change display name
		settingsGroup.PATCH("/account/display-name/", userRouter.ChangeDisplayName)
		// Allow user to change password.
		settingsGroup.PATCH("/account/password/", userRouter.UpdatePassword)

		// Show full account data.
		settingsGroup.GET("/profile/", userRouter.Profile)
	}

	// API access control
	apiRouter := controller.NewOAuthRouter(gormDBs)
	oauthGroup := apiGroup.Group("/oauth", guard.RequireLoggedIn)
	{
		// Get a list of apps. /apps?page=<int>&per_page=<int>
		oauthGroup.GET("/apps/", apiRouter.ListApps)
		// Create a new app
		oauthGroup.POST("/apps/", apiRouter.CreateApp)
		// Get a specific app
		oauthGroup.GET("/apps/:id/", apiRouter.LoadApp)
		// Update an app
		oauthGroup.PATCH("/apps/:id/", apiRouter.UpdateApp)
		// Deactivate an app
		oauthGroup.DELETE("/apps/:id/", apiRouter.RemoveApp)

		// Get a list access tokens.
		// /api/keys?client_id=<string>&page=<number>&per_page=<number>
		oauthGroup.GET("/keys/", apiRouter.ListKeys)
		// Create a new key.
		oauthGroup.POST("/keys/", apiRouter.CreateKey)
		// Delete a single key belong to an app or a human.
		// A key could only be deleted by its owner, regardless of
		// being an app's access token or a personal key.
		oauthGroup.DELETE("/keys/:id/", apiRouter.RemoveKey)
	}

	// A reader's profile.
	readersGroup := apiGroup.Group("/readers", guard.RequireLoggedIn)
	{
		// Search ftc account: /search/reader?q=<email|username|phone>&kind=ftc
		// Search wx account: /search/reader?q=<nickname>&kind=wechat&page=<number>&per_page=<number>
		readersGroup.GET("/search/", readerRouter.SearchAccount)
		// Get an ftc account
		readersGroup.GET("/ftc/:id/", readerRouter.LoadFTCAccount)
		// Get more details of an ftc account
		readersGroup.GET("/ftc/:id/profile/", readerRouter.LoadFtcProfile)

		// Load a wechat Account
		readersGroup.GET("/wx/:id/", readerRouter.LoadWxAccount)
		// Load more details of a wechat account
		readersGroup.GET("/wx/:id/profile/", readerRouter.LoadWxProfile)
	}

	sandboxGroup := apiGroup.Group("/sandbox", guard.RequireLoggedIn)
	{
		sandboxGroup.POST("/", readerRouter.CreateTestUser)
		sandboxGroup.GET("/", readerRouter.ListTestUsers)
		sandboxGroup.GET("/:id/", readerRouter.LoadTestAccount)
		sandboxGroup.DELETE("/:id/", readerRouter.DeleteTestAccount)
		// Change sandbox user password. This is like a force override.
		sandboxGroup.PATCH("/:id/password/", readerRouter.ChangeSandboxPassword)
	}

	// Manipulate membership.
	// The `id` in this section should be ftc id if exists,
	// and fallback to wechat union id if ftc id does not exist.
	// This section cannot use the standard restful way of
	// setting id as an url parameters due to membership has
	// double ids. This is deficiency of the initial design.
	// We could only attach two optional ids to query parameters:
	// ?ftc_id=<string>&union_id=<string>
	memberGroup := apiGroup.Group("/memberships", guard.RequireLoggedIn)
	{
		// Update an ftc subscription or create one if not present.
		// The membership might be under email account or wechat account.
		// Client should pass all ids if so that we could determine how to find out user account.
		memberGroup.POST("/", readerRouter.UpsertFtcMember)
		// Delete a membership.
		// It is assumed you are deleting an FTC member, which will be denied if it is not purchased via ali or wx pay.
		memberGroup.DELETE("/:id/", readerRouter.DeleteFtcMember)
	}

	snapshotGroup := apiGroup.Group("/snapshots", guard.RequireLoggedIn)
	{
		// ?ftc_id=<string>&union_id=<string>&page=<int>&per_page=<int>
		snapshotGroup.GET("/", readerRouter.ListSnapshots, xhttp.RequireUserIDsQuery)
	}

	// `id` for iapGroup is for original transaction id
	iapGroup := apiGroup.Group("/iap", guard.RequireLoggedIn)
	{
		// List IAP.
		// ?page=<int>&per_page=<int>
		// X-User-Id is required.
		iapGroup.GET("/", readerRouter.ListIAPSubs, xhttp.RequireUserIDHeader)
		// Load a single IAP subscription.
		iapGroup.GET("/:id/", readerRouter.LoadIAPSubs)
		// Refresh an existing IAP.
		iapGroup.PATCH("/:id/", readerRouter.RefreshIAPSubs)

		// Link iap to an ftc account.
		iapGroup.POST("/:id/link/", readerRouter.LinkIAP)
		iapGroup.POST("/:id/unlink/", readerRouter.UnlinkIAP)
	}

	// Paywall, products, prices, discounts
	paywallGroup := apiGroup.Group("/paywall", guard.RequireLoggedIn)
	{
		paywallGroup.GET("/", productRoutes.LoadPaywall)

		// Create a banner
		paywallGroup.POST("/banner/", productRoutes.CreateBanner)
		paywallGroup.POST("/banner/promo/", productRoutes.CreatePromoBanner)
		// Drop promo from a banner
		paywallGroup.DELETE("/banner/promo/", productRoutes.DropPromoBanner)

		// Create, list, update products.
		// All path have query parameter `?live=<true|false>`. Default true.
		productGroup := paywallGroup.Group("/products")
		{
			// Create a product
			productGroup.POST("/", productRoutes.CreateProduct)
			// List all products. The product has a plan field. The plan does not contain discount.
			productGroup.GET("/", productRoutes.ListProducts)
			// Retrieve a product by id.
			productGroup.GET("/:productId/", productRoutes.LoadProduct)
			// Put a product on paywall.
			productGroup.POST("/:productId/activate/", productRoutes.ActivateProduct)
			// Update a product.
			productGroup.PATCH("/:productId/", productRoutes.UpdateProduct)
			// Attached an introductory price to a product
			// Deprecated. This may no longer works
			productGroup.PATCH("/:productId/intro/", productRoutes.AttachIntroPrice)
			// Delete an introductory price of a product.
			// Deprecated. This may no longer works.
			productGroup.DELETE("/:productId/intro/", productRoutes.DropIntroPrice)
		}

		// Create, list plans and its discount.
		priceGroup := paywallGroup.Group("/prices")
		{
			// Create a price for a product
			priceGroup.POST("/", productRoutes.CreatePrice)
			// List all prices under a product.
			// ?product_id=<string>&live=<true|false>
			priceGroup.GET("/", productRoutes.ListPriceOfProduct)

			priceGroup.GET("/:priceId/", productRoutes.LoadPrice)
			// Turn a price into active state under a product.
			// There's only one edition of active price under a specific product.
			priceGroup.POST("/:priceId/activate/", productRoutes.ActivatePrice)
			priceGroup.POST("/:priceId/deactivate/", productRoutes.DeactivatePrice)
			priceGroup.PATCH("/:priceId/", productRoutes.UpdatePrice)
			priceGroup.PATCH("/:priceId/discounts/", productRoutes.RefreshPriceDiscounts)
			priceGroup.DELETE("/:priceId/", productRoutes.ArchivePrice)
		}

		discountGroup := paywallGroup.Group("/discounts")
		{
			discountGroup.POST("/", productRoutes.CreateDiscount)
			discountGroup.DELETE("/:id/", productRoutes.RemoveDiscount)
		}
	}

	stripeGroup := apiGroup.Group("/stripe", guard.RequireLoggedIn)
	{
		// ?page=<int>&per_page=<int>&live=<bool>
		stripeGroup.GET("/prices/", productRoutes.ListStripePrices)

		// TODO: create a stripe price directly here
		// by sending data to Stripe API rather thatn
		// creating one in Stripe dashboard.
		// stripeGroup.PUT("/prices/", productRoutes.CreateStripePrice)

		// ?live=<bool>&refersh=<bool>
		stripeGroup.GET("/prices/:id/", productRoutes.LoadStripePrice)
		// Update stripe price metadata.
		// ?live=<bool>
		stripeGroup.PATCH("/prices/:id/", productRoutes.UpdateStripePriceMeta)
		// Activate a stripe price.
		// ?live=<bool>
		stripeGroup.PATCH("/prices/:id/activate/", productRoutes.ActivateStripePrice)
		// Deactivate a stripe price.
		// ?live=<bool>
		stripeGroup.PATCH("/prices/:id/deactivate/", productRoutes.DeactivateStripePrice)
		// List coupons under a price.
		// ?live=<bool>
		stripeGroup.GET("/prices/:id/coupons/", productRoutes.ListStripeCoupons)

		// Load a stripe coupon
		// ?live=<bool>&refresh=<bool>
		stripeGroup.GET("/coupons/:id/", productRoutes.LoadStripeCoupon)
		// Update a coupon
		// ?live=<bool>
		stripeGroup.POST("/coupons/:id/", productRoutes.UpdateCoupon)
		// Activate a coupon
		// ?live=<bool>
		stripeGroup.PATCH("/coupons/:id/activate/", productRoutes.ActivateCoupon)
		// Delete a coupon
		// ?live=<bool>
		stripeGroup.DELETE("/coupons/:id/", productRoutes.DeleteCoupon)
	}

	androidRouter := controller.NewAndroidRouter(
		apiClients.Select(true),
		logger)
	androidGroup := apiGroup.Group("/android", guard.RequireLoggedIn)
	{
		androidGroup.POST("/releases/", androidRouter.CreateRelease)
		androidGroup.GET("/releases/", androidRouter.ListReleases)
		androidGroup.GET("/releases/:versionName/", androidRouter.ReleaseOf)
		androidGroup.PATCH("/releases/:versionName/", androidRouter.UpdateRelease)
		androidGroup.DELETE("/releases/:versionName/", androidRouter.DeleteRelease)
	}

	wikiRouter := controller.NewWikiRouter(gormDBs)
	wikiGroup := apiGroup.Group("/wiki", guard.RequireLoggedIn)
	{
		wikiGroup.GET("/", wikiRouter.ListArticle)
		wikiGroup.POST("/", wikiRouter.CreateArticle)
		wikiGroup.GET("/:id/", wikiRouter.OneArticle)
		wikiGroup.PATCH("/:id/", wikiRouter.UpdateArticle)
	}

	legalRouter := controller.NewLegalRoutes(
		apiClients.Select(true),
		logger)
	legalGroup := apiGroup.Group("/legal", guard.RequireLoggedIn)
	{
		// ?page=<int>&per_page=<int>
		legalGroup.GET("/", legalRouter.List)
		legalGroup.POST("/", legalRouter.Create)
		legalGroup.GET("/:id/", legalRouter.Load)
		legalGroup.PATCH("/:id/", legalRouter.Update)
		legalGroup.POST("/:id/publish/", legalRouter.Publish)
	}

	e.Logger.Fatal(e.Start(":3001"))
}

// RestfulErrorHandler implements echo's HTTPErrorHandler.
func errorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var re *render.ResponseError
	switch {
	case errors.As(err, &re):
		if re.Message == "" {
			re.Message = http.StatusText(re.StatusCode)
		}

	default:
		re = render.NewInternalError(err.Error())
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(re.StatusCode)
	} else {
		err = c.JSON(re.StatusCode, re)
	}
	if err != nil {
		c.Logger().Error(err)
	}
}
