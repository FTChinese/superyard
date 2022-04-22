package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
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
	"net/http"
	"os"
	"time"
)

//go:embed build/api.toml
var tomlConfig string

//go:embed client_version_next
var clientVersionNext string

var clientVersionNg string

var (
	isProduction bool
	version      string
	build        string
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

	myDB := db.MustNewMyDBs(isProduction)

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
		Repo:       readers.New(myDB, logger),
		Postman:    hanqiPm,
		APIClients: apiClients,
		Logger:     logger,
	}
	productRoutes := controller.NewPaywallRouter(apiClients, logger)

	userRouter := controller.NewUserRouter(myDB, ftcPm, guard)

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

	// Staff administration
	adminRouter := controller.NewAdminRouter(myDB, ftcPm)
	adminGroup := apiGroup.Group("/admin", guard.RequireLoggedIn)
	{
		//	GET /staff?page=<number>&per_page=<number>
		adminGroup.GET("/staff/", adminRouter.ListStaff)
		// Create a staff
		adminGroup.POST("/staff/", adminRouter.CreateStaff)

		// Get the staff profile
		adminGroup.GET("/staff/:id/", adminRouter.StaffProfile)
		// UpdateProfile a staff's profile
		adminGroup.PATCH("/staff/:id/", adminRouter.UpdateStaff)
		// Delete a staff.
		adminGroup.DELETE("/staff/:id/", adminRouter.DeleteStaff)
		// Reinstate a deactivated staff
		adminGroup.PUT("/staff/:id/", adminRouter.Reinstate)

		adminGroup.GET("/vip/", adminRouter.ListVIPs)
		adminGroup.PUT("/vip/:id/", adminRouter.SetVIP(true))
		adminGroup.DELETE("/vip/:id/", adminRouter.SetVIP(false))
	}

	// API access control
	apiRouter := controller.NewOAuthRouter(myDB)
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
		// ?q=<email|username>
		readersGroup.GET("/ftc/", readerRouter.FindFTCAccount)
		readersGroup.GET("/ftc/:id/", readerRouter.LoadFTCAccount)
		readersGroup.GET("/ftc/:id/profile/", readerRouter.LoadFtcProfile)
		// Login history
		readersGroup.GET("/ftc/:id/activities/", readerRouter.LoadActivities)

		// Wx Account
		readersGroup.GET("/wx/:id/", readerRouter.LoadWxAccount)
		readersGroup.GET("/wx/:id/profile/", readerRouter.LoadWxProfile)
		// Wx login history
		readersGroup.GET("/wx/:id/login/", readerRouter.LoadOAuthHistory)
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
		memberGroup.POST("/", readerRouter.CreateFtcMember)
		memberGroup.PATCH("/:id/", readerRouter.UpdateFtcMember)
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

		// The membership linked to an original transaction id.
		iapGroup.GET("/:id/membership/", readerRouter.IAPMember)

		// Link iap to an ftc account.
		iapGroup.POST("/:id/link/", readerRouter.LinkIAP)
		iapGroup.POST("/:id/unlink/", readerRouter.UnlinkIAP)
	}

	orderGroup := apiGroup.Group("/orders", guard.RequireLoggedIn)
	{
		// Get a list of orders of a specific reader.
		// ?ftc_id=<string>&union_id=<string>&page=<int>&per_page=<int>
		// ftc_id and union_id are not both required,
		// but at least one should be present.
		orderGroup.GET("/", readerRouter.ListOrders)

		// Get an order
		// This can also be used to search an order by id.
		orderGroup.GET("/:id/", readerRouter.LoadOrder)
		orderGroup.GET("/:id/webhook/alipay/", readerRouter.AliWebhook)
		orderGroup.GET("/:id/webhook/wechat/", readerRouter.WxWebhook)
		// Confirm an order. This also renew or upgrade membership.
		orderGroup.PATCH("/:id/", readerRouter.ConfirmOrder)
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
			productGroup.POST("/:productId/", productRoutes.ActivateProduct)
			// Update a product.
			productGroup.PATCH("/:productId/", productRoutes.UpdateProduct)
			// Attached an introductory price to a product
			productGroup.PATCH("/:productId/intro/", productRoutes.AttachIntroPrice)
			// Delete an introductory price of a product.
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
			// Turn a price into active state under a product.
			// There's only one edition of active price under a specific product.
			priceGroup.POST("/:priceId/", productRoutes.ActivatePrice)
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
		stripeGroup.GET("/prices/", productRoutes.ListStripePrices)
		stripeGroup.GET("/prices/:id/", productRoutes.LoadStripePrice)
		stripeGroup.GET("/prices/:id/coupons/", productRoutes.ListStripeCoupons)

		stripeGroup.GET("/coupons/:id/", productRoutes.LoadStripeCoupon)
		stripeGroup.POST("/coupons/:id/", productRoutes.UpdateCoupon)
		stripeGroup.DELETE("/coupons/:id/", productRoutes.DeleteCoupon)
	}

	b2bRouter := controller.NewB2BRouter(isProduction)
	b2bGroup := apiGroup.Group("/b2b", guard.RequireLoggedIn)
	{
		b2bGroup.GET("/teams/:id/", b2bRouter.LoadTeam)
		b2bGroup.GET("/orders/", b2bRouter.ListOrders)
		b2bGroup.GET("/orders/:id/", b2bRouter.LoadOrder)
		b2bGroup.POST("/orders/:id/", b2bRouter.ConfirmOrder)
	}

	androidRouter := controller.NewAndroidRouter(myDB)
	androidGroup := apiGroup.Group("/android", guard.RequireLoggedIn)
	{
		androidGroup.GET("/gh/latest/", androidRouter.GHLatestRelease)
		androidGroup.GET("/gh/tags/:tag/", androidRouter.GHRelease)

		androidGroup.GET("/exists/:versionName/", androidRouter.TagExists)
		androidGroup.POST("/releases/", androidRouter.CreateRelease)
		androidGroup.GET("/releases/", androidRouter.Releases)
		androidGroup.GET("/releases/:versionName/", androidRouter.SingleRelease)
		androidGroup.PATCH("/releases/:versionName/", androidRouter.UpdateRelease)
		androidGroup.DELETE("/releases/:versionName/", androidRouter.DeleteRelease)
	}

	wikiRouter := controller.NewWikiRouter(myDB)
	wikiGroup := apiGroup.Group("/wiki", guard.RequireLoggedIn)
	{
		wikiGroup.GET("/", wikiRouter.ListArticle)
		wikiGroup.POST("/", wikiRouter.CreateArticle)
		wikiGroup.GET("/:id/", wikiRouter.OneArticle)
		wikiGroup.PATCH("/:id/", wikiRouter.UpdateArticle)
	}

	statsRouter := controller.NewStatsRouter(myDB)
	statsGroup := apiGroup.Group("/stats")
	{
		statsGroup.GET("/signup/daily/", statsRouter.DailySignUp)
		statsGroup.GET("/income/year/:year/", statsRouter.YearlyIncome)
	}

	whGroup := apiGroup.Group("/webhook", guard.RequireLoggedIn)
	{
		whGroup.GET("/failure/alipay/", statsRouter.AliUnconfirmed)
		whGroup.GET("/failure/wechat/", statsRouter.WxUnconfirmed)
	}

	// Search
	searchGroup := apiGroup.Group("/search")
	{
		// Search by cms user's name: /search/staff?q=<user_name>
		searchGroup.GET("/staff/", adminRouter.Search)
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
