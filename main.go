package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/go-rest/render"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/ftchinese/superyard/pkg/config"
	db2 "gitlab.com/ftchinese/superyard/pkg/db"
	"gitlab.com/ftchinese/superyard/web/views"
	"net/http"
	"os"

	"github.com/FTChinese/go-rest/postoffice"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/superyard/controller"
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

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

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

var mime = map[string]string{
	".js":  "text/javascript",
	".css": "text/css",
}

func main() {

	db := db2.MustNewDB(cfg.MustGetDBConn("mysql.master"))
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

	baseGroup := e.Group("/api")

	userRouter := controller.NewUserRouter(db, post, guard)
	// Login
	// Input {userName: string, password: string}
	baseGroup.POST("/login/", userRouter.Login)
	// Password reset
	baseGroup.POST("/password-reset/", userRouter.ResetPassword)
	baseGroup.POST("/password-reset/letter/", userRouter.ForgotPassword)
	baseGroup.GET("/password-reset/tokens/:token/", userRouter.VerifyResetToken)

	settingsGroup := baseGroup.Group("/settings", guard.RequireLoggedIn)
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
	staffRouter := controller.NewStaffRouter(db, post)
	staffGroup := baseGroup.Group("/staff", guard.RequireLoggedIn)
	{
		//	GET /staff?page=<number>&per_page=<number>
		staffGroup.GET("/", staffRouter.List)
		// Create a staff
		staffGroup.POST("/", staffRouter.Create)

		// Get the staff profile
		staffGroup.GET("/:id/", staffRouter.Profile)
		// UpdateProfile a staff's profile
		staffGroup.PATCH("/:id/", staffRouter.Update)
		// Delete a staff.
		staffGroup.DELETE("/:id/", staffRouter.Delete)
		// Reinstate a deactivated staff
		staffGroup.PUT("/:id/", staffRouter.Reinstate)
	}

	// API access control
	apiRouter := controller.NewOAuthRouter(db)
	oauthGroup := baseGroup.Group("/oauth", guard.RequireLoggedIn)
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

	readerRouter := controller.NewReaderRouter(db, hanqi)
	// A reader's profile.
	readersGroup := baseGroup.Group("/readers", guard.RequireLoggedIn)
	{
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

	memberRouter := controller.NewMemberRouter(db)
	memberGroup := baseGroup.Group("/memberships", guard.RequireLoggedIn)
	{
		// Create a new membership:
		// Input: {ftcId: string,
		// unionId: string,
		// tier: string,
		// cycle: string,
		// expireDate: string,
		// payMethod: string
		// stripeSubId: string,
		// stripePlanId: string,
		// autoRenewal: boolean,
		// status: ""}
		memberGroup.POST("/", memberRouter.CreateMember)
		// Get one subscription
		memberGroup.GET("/:id/", memberRouter.LoadMember)
		// UpdateProfile a subscription
		memberGroup.PATCH("/:id/", memberRouter.UpdateMember)
		// Delete a subscription
		memberGroup.DELETE("/:id/", memberRouter.DeleteMember)
	}

	orderGroup := baseGroup.Group("/orders", guard.RequireLoggedIn)
	{
		// Get a list of orders of a specific reader.
		// /orders?ftc_id=<string>&union_id=<string>&page=<int>&per_page=<int>
		// ftc_id and union_id are not both required,
		// but at least one should be present.
		baseGroup.GET("/", readerRouter.ListOrders)

		// Get an order
		// This can also be used to search an order by id.
		orderGroup.GET("/:id/", readerRouter.LoadOrder)
		// Confirm an order. This also renew or upgrade
		// membership.
		orderGroup.PATCH("/:id/", readerRouter.ConfirmOrder)
	}

	androidRouter := controller.NewAndroidRouter(db)
	androidGroup := baseGroup.Group("/android", guard.RequireLoggedIn)
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

	wikiRouter := controller.NewWikiRouter(db)
	wikiGroup := baseGroup.Group("/wiki", guard.RequireLoggedIn)
	{
		wikiGroup.GET("/", wikiRouter.ListArticle)
		wikiGroup.POST("/", wikiRouter.CreateArticle)
		wikiGroup.GET("/:id", wikiRouter.OneArticle)
		wikiGroup.PATCH("/:id", wikiRouter.UpdateArticle)
	}

	statsRouter := controller.NewStatsRouter(db)
	statsGroup := baseGroup.Group("/stats")
	{
		statsGroup.GET("/signup/daily/", statsRouter.DailySignUp)
		statsGroup.GET("/income/year/{year}/", statsRouter.YearlyIncome)
	}

	// Search
	searchGroup := baseGroup.Group("/search")
	{
		// Search by cms user's name: /search/staff?q=<user_name>
		searchGroup.GET("/staff/", staffRouter.Search)
		// Search ftc account: /search/reader?q=<email>&kind=ftc
		// Search wx account: /search/reader?q=<nickname>&kind=wechat&page=<number>&per_page=<number>
		searchGroup.GET("/reader/", readerRouter.SearchAccount)
	}

	e.Logger.Fatal(e.Start(":3100"))
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
