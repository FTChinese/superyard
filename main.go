package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/go-rest/render"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/FTChinese/go-rest/postoffice"
	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/superyard/controller"
	"gitlab.com/ftchinese/superyard/models/util"
)

var (
	isProduction bool
	version      string
	build        string
	config       Config
	logger       = logrus.WithField("project", "superyard").WithField("package", "main")
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

	config = Config{
		Debug:   !isProduction,
		Version: version,
		BuiltAt: build,
		Year:    0,
	}

	controller.HomeData.Debug = !isProduction
}

func main() {

	db, err := util.NewDBX(config.MustGetDBConn("mysql.master"))
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	emailConn := MustGetEmailConn()
	post := postoffice.NewPostman(
		emailConn.Host,
		emailConn.Port,
		emailConn.User,
		emailConn.Pass)

	e := echo.New()
	e.HTTPErrorHandler = errorHandler

	if !isProduction {
		e.Static("/", "build/dev")
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CSRF())

	e.GET("/", controller.Home)

	baseGroup := e.Group("/api")

	userRouter := controller.NewUserRouter(db, post)
	// Login
	// Input {userName: string, password: string}
	baseGroup.POST("/login", userRouter.Login)
	// Password reset
	baseGroup.POST("/password-reset", userRouter.ResetPassword)
	baseGroup.POST("/password-reset/letter", userRouter.ForgotPassword)
	baseGroup.GET("/password-reset/tokens/:token", userRouter.VerifyToken)

	settingsGroup := baseGroup.Group("/settings", controller.CheckJWT)

	// Use to renew Json Web Token
	settingsGroup.GET("/account", userRouter.Account)
	// Set email if empty. User can only set
	// it once.
	settingsGroup.PATCH("/account/email", userRouter.SetEmail)
	// Allow user to change display name
	settingsGroup.PATCH("/account/display-name", userRouter.ChangeDisplayName)
	// Allow user to change password.
	settingsGroup.PATCH("/account/password", userRouter.ChangePassword)

	// Show full account data.
	settingsGroup.GET("/profile", userRouter.Profile)

	// Staff administration
	staffRouter := controller.NewStaffRouter(db, post)
	// User data.
	staffGroup := baseGroup.Group("/staff")
	//	GET /staff?page=<number>&per_page=<number>
	staffGroup.GET("/", staffRouter.List)
	// Create a staff
	staffGroup.POST("/", staffRouter.Create)
	// Get the staff profile
	staffGroup.GET("/:id", staffRouter.Profile)
	// UpdateProfile a staff's profile
	staffGroup.PATCH("/:id", staffRouter.Update)
	// Delete a staff.
	staffGroup.DELETE("/:id", staffRouter.Delete)
	// Reinstate a deactivated staff
	staffGroup.PUT("/:id", staffRouter.Reinstate)

	// API access control
	apiRouter := controller.NewOAuthRouter(db)

	oauthGroup := baseGroup.Group("/oauth")

	// Get a list of apps. /apps?page=<int>&per_page=<int>
	oauthGroup.GET("/apps", apiRouter.ListApps)
	// Create a new app
	oauthGroup.POST("/apps", apiRouter.CreateApp)
	// Get a specific app
	oauthGroup.GET("/apps/:id", apiRouter.LoadApp)
	// Update an app
	oauthGroup.PATCH("/apps/:id", apiRouter.UpdateApp)
	// Deactivate an app
	oauthGroup.DELETE("/apps/:id", apiRouter.RemoveApp)

	// Get a list access tokens.
	// /api/keys?client_id=<string>&page=<number>&per_page=<number>
	// /api/keys?staff_name=<string>&page=<number>&per_page=<number>
	oauthGroup.GET("/keys", apiRouter.ListKeys)
	// Create a new key.
	oauthGroup.POST("/keys", apiRouter.CreateKey)
	// Delete a single key belong to an app or a human.
	// A key could only be deleted by its owner, regardless of
	// being an app's access token or a personal key.
	oauthGroup.DELETE("/keys/:id", apiRouter.RemoveKey)

	readerRouter := controller.NewReaderRouter(db)
	// Handle VIPs
	vipGroup := baseGroup.Group("/vip")
	vipGroup.GET("/", readerRouter.ListVIP)
	vipGroup.PUT("/:id", readerRouter.GrantVIP)
	vipGroup.DELETE("/:id", readerRouter.RevokeVIP)

	// A reader's profile.
	readersGroup := baseGroup.Group("/readers")

	readersGroup.GET("/ftc/:id", readerRouter.LoadFTCAccount)
	readersGroup.GET("/ftc/:id/profile", readerRouter.LoadFtcProfile)
	// Login history
	readersGroup.GET("/ftc/:id/activities", readerRouter.LoadActivities)

	// Wx Account
	readersGroup.GET("/wx/:id", readerRouter.LoadWxAccount)
	readersGroup.GET("/wx/:id/profile", readerRouter.LoadWxProfile)
	// Wx login history
	readersGroup.GET("/wx/:id/login", readerRouter.LoadOAuthHistory)

	memberRouter := controller.NewMemberRouter(db)
	memberGroup := baseGroup.Group("/memberships")
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
	memberGroup.GET("/:id", memberRouter.LoadMember)
	// UpdateProfile a subscription
	memberGroup.PATCH("/:id", memberRouter.UpdateMember)
	// Delete a subscription
	memberGroup.DELETE("/:id", memberRouter.DeleteMember)

	orderRouter := controller.NewOrderRouter(db)
	orderGroup := baseGroup.Group("/orders")
	// Get a list of orders of a specific reader.
	// /orders?ftc_id=<string>&union_id=<string>&page=<int>&per_page=<int>
	// ftc_id and union_id are not both required,
	// but at least one should be present.
	orderGroup.GET("/", orderRouter.ListOrders)
	// Create an order
	orderGroup.POST("/", orderRouter.CreateOrder)
	// Get an order
	orderGroup.GET("/:id", orderRouter.LoadOrder)
	// Confirm an order. This also renew or upgrade
	// membership.
	orderGroup.PATCH("/:id", orderRouter.ConfirmOrder)

	promoRouter := controller.NewPromoRouter(db)
	promoGroup := baseGroup.Group("/promos")
	// ListStaff promos by page
	promoGroup.GET("/", promoRouter.ListPromos)
	// Create a new promo
	promoGroup.POST("/", promoRouter.CreateSchedule)
	// Get a promo
	promoGroup.GET("/:id", promoRouter.LoadPromo)
	// Delete a promo
	promoGroup.DELETE("/:id", promoRouter.DisablePromo)
	promoGroup.PATCH("/:id/plans", promoRouter.SetPricingPlans)
	promoGroup.PATCH("/:id/banner", promoRouter.SetBanner)

	androidRouter := controller.NewAndroidRouter(db)
	androidGroup := baseGroup.Group("/android")

	androidGroup.GET("/gh/latest", androidRouter.GHLatestRelease)
	androidGroup.GET("/gh/tags/:tag", androidRouter.GHRelease)

	androidGroup.GET("/exists/:versionName", androidRouter.TagExists)
	androidGroup.POST("/releases", androidRouter.CreateRelease)
	androidGroup.GET("/releases", androidRouter.Releases)
	androidGroup.GET("/releases/:versionName", androidRouter.SingleRelease)
	androidGroup.PATCH("/releases/:versionName", androidRouter.UpdateRelease)
	androidGroup.DELETE("/releases/:versionName", androidRouter.DeleteRelease)

	statsRouter := controller.NewStatsRouter(db)
	statsGroup := baseGroup.Group("/stats")
	statsGroup.GET("/signup/daily", statsRouter.DailySignUp)
	statsGroup.GET("/income/year/{year}", statsRouter.YearlyIncome)

	// Search
	searchGroup := baseGroup.Group("/search")
	// Search by cms user's name: /search/staff?name=<user_name>
	searchGroup.GET("/staff", staffRouter.Search)
	// Search ftc account: /search/reader?q=<email>&kind=ftc
	// Search wx account: /search/reader?q=<nickname>&kind=wechat&page=<number>&per_page=<number>
	searchGroup.GET("/reader", readerRouter.SearchAccount)

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
