package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/go-rest/render"
	"github.com/FTChinese/superyard/internal/repository/subsapi"
	"github.com/FTChinese/superyard/pkg/config"
	"github.com/FTChinese/superyard/pkg/db"
	"github.com/FTChinese/superyard/web/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"

	"github.com/FTChinese/go-rest/postoffice"
	"github.com/spf13/viper"

	"github.com/FTChinese/superyard/internal/controller"
)

var (
	isProduction bool
	version      string
	build        string
	cfg          config.Config
)

func init() {
	flag.BoolVar(&isProduction, "production", false, "Indicate productions environment if present")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		os.Exit(1)
	}

	cfg = config.Config{
		Debug:   !isProduction,
		Version: version,
		BuiltAt: build,
		Year:    0,
	}
}

func main() {

	sqlDB := db.MustNewDB(cfg.MustGetDBConn("mysql.master"))
	post := postoffice.New(config.MustGetEmailConn())
	hanqi := postoffice.New(config.MustGetHanqiConn())

	guard := controller.MustNewGuard()

	e := echo.New()
	e.Renderer = views.New()

	if !isProduction {
		e.Static("/static", "build/public/static")
	}

	e.Pre(middleware.AddTrailingSlash())
	e.HTTPErrorHandler = errorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CSRF())
	e.Use(controller.DumpRequest)

	e.GET("/*", controller.Home)

	apiGroup := e.Group("/api")

	subsAPI := subsapi.NewClient(cfg.Debug)

	userRouter := controller.NewUserRouter(sqlDB, post, guard)
	// Login
	// Input {userName: string, password: string}
	apiGroup.POST("/login/", userRouter.Login)
	// Password reset
	apiGroup.POST("/password-reset/", userRouter.ResetPassword)
	apiGroup.POST("/password-reset/letter/", userRouter.ForgotPassword)
	apiGroup.GET("/password-reset/tokens/:token/", userRouter.VerifyResetToken)

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
	adminRouter := controller.NewAdminRouter(sqlDB, post)
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
	apiRouter := controller.NewOAuthRouter(sqlDB)
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
		// /api/keys?staff_name=<string>&page=<number>&per_page=<number>
		oauthGroup.GET("/keys/", apiRouter.ListKeys)
		// Create a new key.
		oauthGroup.POST("/keys/", apiRouter.CreateKey)
		// Delete a single key belong to an app or a human.
		// A key could only be deleted by its owner, regardless of
		// being an app's access token or a personal key.
		oauthGroup.DELETE("/keys/:id/", apiRouter.RemoveKey)
	}

	readerRouter := controller.NewReaderRouter(sqlDB, hanqi, subsAPI)
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

	// The id in this section should be ftc id if exists in user account, and then
	// use wechat union id if ftc id does not exist.
	memberGroup := apiGroup.Group("/memberships", guard.RequireLoggedIn)
	{
		// Update a ftc subscription or create one if not present.
		// The membership might be under email account or wechat account.
		// Client should pass all ids if so that we could determine how to find out user account.
		memberGroup.POST("/", readerRouter.UpsertFtcSubs)
		// Get a reader's membership by compound id.
		memberGroup.GET("/:id/", readerRouter.LoadMember)
		// Delete the sandbox user membership, not matter what it is.
		memberGroup.DELETE("/:id/", readerRouter.DeleteMember)

		// Link user to IAP.
		memberGroup.PATCH("/:id/apple/", readerRouter.LinkIAP)
		// Add stripe subscription it.
		memberGroup.PATCH("/:id/stripe/", readerRouter.UpsertStripeSubs)
	}

	iapGroup := apiGroup.Group("/iap", guard.RequireLoggedIn)
	{
		// List IAP.
		// ?page=<int>&per_page=<int>
		iapGroup.GET("/", readerRouter.ListIAPSubs)
		// There is not POST for IAP since you cannot create one here.
		// Load a single IAP
		iapGroup.GET("/:id/", readerRouter.LoadIAPSubs)
		// Refresh an existing IAP.
		iapGroup.PATCH("/:id/", readerRouter.RefreshIAPSubs)
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

	orderGroup := apiGroup.Group("/orders", guard.RequireLoggedIn)
	{
		// Get a list of orders of a specific reader.
		// /orders?ftc_id=<string>&union_id=<string>&page=<int>&per_page=<int>
		// ftc_id and union_id are not both required,
		// but at least one should be present.
		apiGroup.GET("/", readerRouter.ListOrders)

		// Get an order
		// This can also be used to search an order by id.
		orderGroup.GET("/:id/", readerRouter.LoadOrder)
		// Confirm an order. This also renew or upgrade
		// membership.
		orderGroup.PATCH("/:id/", readerRouter.ConfirmOrder)
	}

	productRouter := controller.NewProductRouter(sqlDB, subsAPI)
	paywallGroup := apiGroup.Group("/paywall", guard.RequireLoggedIn)
	{
		paywallGroup.GET("/", productRouter.LoadPaywall)

		// Requesting subscription api to bust cached paywall data.
		paywallGroup.GET("/build/", productRouter.RefreshAPI)

		// Create a banner
		paywallGroup.POST("/banner/", productRouter.CreateBanner)
		// Retrieve a banner
		paywallGroup.GET("/banner/", productRouter.LoadBanner)
		// Update a banner
		paywallGroup.PATCH("/banner/", productRouter.UpdateBanner)
		// Drop promo from a banner
		paywallGroup.DELETE("/banner/promo/", productRouter.DropBannerPromo)

		// Create a promo
		paywallGroup.POST("/promo/", productRouter.CreatePromo)
		// Load a promo
		paywallGroup.GET("/promo/:id/", productRouter.LoadPromo)

		// A list of active plans shown on paywall.
		paywallGroup.GET("/plans/", productRouter.ListPlansOnPaywall)
	}

	// Create, list, update products.
	productGroup := apiGroup.Group("/products", guard.RequireLoggedIn)
	{
		// Create a product
		productGroup.POST("/", productRouter.CreateProduct)
		// List all products. The product has a plan field. The plan does not contains discount.
		productGroup.GET("/", productRouter.ListPricedProducts)
		// Retrieve a product by id.
		productGroup.GET("/:productId/", productRouter.LoadProduct)
		// Update a product.
		productGroup.PATCH("/:productId/", productRouter.UpdateProduct)
		// Put a product on paywall.
		productGroup.PUT("/:productId/", productRouter.ActivateProduct)
		productGroup.PUT("/:productId/", productRouter.ActivateProduct)
	}

	// Create, list plans and its discount.
	planGroup := apiGroup.Group("/plans", guard.RequireLoggedIn)
	{
		// Create a plan for a product
		planGroup.POST("/", productRouter.CreatePlan)
		// List all plans under a product.
		// ?product_id=<string>
		planGroup.GET("/", productRouter.ListPlansOfProduct)
		// TODO: update a plan's description
		// Set a plan a default one so that it is visible on paywall.
		planGroup.PUT("/:planId/", productRouter.ActivatePlan)

		// Create a discount for a plan and apply to it.
		planGroup.POST("/:planId/discount/", productRouter.CreateDiscount)
		planGroup.DELETE("/:planId/discount/", productRouter.DropDiscount)
	}

	androidRouter := controller.NewAndroidRouter(sqlDB)
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

	wikiRouter := controller.NewWikiRouter(sqlDB)
	wikiGroup := apiGroup.Group("/wiki", guard.RequireLoggedIn)
	{
		wikiGroup.GET("/", wikiRouter.ListArticle)
		wikiGroup.POST("/", wikiRouter.CreateArticle)
		wikiGroup.GET("/:id/", wikiRouter.OneArticle)
		wikiGroup.PATCH("/:id/", wikiRouter.UpdateArticle)
	}

	statsRouter := controller.NewStatsRouter(sqlDB)
	statsGroup := apiGroup.Group("/stats")
	{
		statsGroup.GET("/signup/daily/", statsRouter.DailySignUp)
		statsGroup.GET("/income/year/:year/", statsRouter.YearlyIncome)
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
	re, ok := err.(*render.ResponseError)
	if !ok {
		re = render.NewInternalError(err.Error())
	}

	if re.Message == "" {
		re.Message = http.StatusText(re.StatusCode)
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(re.StatusCode)
		} else {
			err = c.JSON(re.StatusCode, re)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}
}
