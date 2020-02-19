package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"

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
}

func main() {

	db, err := util.NewDBX(config.MustGetDBConn("mysql.master"))
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	//apnDB, err := util.NewDBX(config.MustGetDBConn("mysql.apn"))
	//if err != nil {
	//	logger.Error(err)
	//	os.Exit(1)
	//}

	emailConn := MustGetEmailConn()
	post := postoffice.NewPostman(
		emailConn.Host,
		emailConn.Port,
		emailConn.User,
		emailConn.Pass)

	e := echo.New()
	e.Renderer = MustNewRenderer(config)
	e.HTTPErrorHandler = util.RestfulErrorHandler

	if !isProduction {
		e.Static("/", "build/dev")
	}

	e.Use(middleware.Logger())
	e.Use(session.Middleware(
		sessions.NewCookieStore(
			[]byte(MustGetSessionKey()),
		),
	))
	e.Use(middleware.Recover())
	//e.Use(middleware.CSRF())

	e.GET("/", func(context echo.Context) error {
		return context.Render(http.StatusOK, "base.html", nil)
	})

	apiBase := e.Group("/api")
	apiBase.GET("/", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	staffRouter := controller.NewStaffRouter(db, post)

	// Login
	// Input {userName: string, password: string}
	apiBase.POST("/login", staffRouter.Login)
	// Password reset
	pwGroup := apiBase.Group("/password-reset")
	pwGroup.POST("/", staffRouter.ResetPassword)
	pwGroup.POST("/letter", staffRouter.ForgotPassword)
	pwGroup.GET("/tokens/:token", staffRouter.VerifyToken)

	// User data.
	staffGroup := apiBase.Group("/staff")
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
	staffGroup.PATCH("/:id/password", staffRouter.UpdatePassword)

	// API access control
	apiRouter := controller.APIRouter(db)

	oauthGroup := apiBase.Group("/oauth")

	oauthGroup.GET("/apps", apiRouter.ListApps)
	oauthGroup.POST("/apps", apiRouter.CreateApp)
	oauthGroup.GET("/apps/{id}", apiRouter.LoadApp)
	oauthGroup.PATCH("/apps/{id}", apiRouter.UpdateApp)
	oauthGroup.DELETE("/apps/{id}", apiRouter.RemoveApp)

	// /api/keys?client_id=<string>&page=<number>&per_page=<number>
	// /api/keys?staff_name=<string>&page=<number>&per_page=<number>
	oauthGroup.GET("/keys", apiRouter.ListKeys)
	// Create a new key.
	oauthGroup.POST("/keys", apiRouter.CreateKey)
	// Delete all keys owned by someone.
	// You cannot delete all keys belonging to an app
	// here since it is performed when an app is deleted.
	oauthGroup.DELETE("/keys", apiRouter.DeletePersonalKeys)
	// Delete a single key belong to an app or a human
	oauthGroup.DELETE("/keys/{id}", apiRouter.RemoveKey)

	//mux := chi.NewRouter()
	//mux.Use(middleware.Logger)
	//mux.Use(middleware.Recoverer)
	//mux.Use(middleware.NoCache)

	readerRouter := controller.NewReaderRouter(db)
	// Handle VIPs
	vipGroup := apiBase.Group("/vip")
	vipGroup.GET("/", readerRouter.ListVIP)
	vipGroup.PUT("/{id}", readerRouter.GrantVIP)
	vipGroup.DELETE("/{id}", readerRouter.RevokeVIP)

	// A reader's profile.
	readersGroup := apiBase.Group("/readers")

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
	memberGroup := apiBase.Group("/memberships")
	// Ge a list of memberships
	memberGroup.GET("/", memberRouter.ListMembers)
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
	memberGroup.GET("/{id}", memberRouter.LoadMember)
	// UpdateProfile a subscription
	memberGroup.PATCH("/{id}", memberRouter.UpdateMember)
	// Delete a subscription
	memberGroup.DELETE("/{id}", memberRouter.DeleteMember)

	orderRouter := controller.NewOrderRouter(db)
	orderGroup := apiBase.Group("/orders")
	// Get a list of orders of a specific reader.
	// /orders?ftc_id=<string>&union_id=<string>&page=<int>&per_page=<int>
	// ftc_id and union_id are not both required,
	// but at least one should be present.
	orderGroup.GET("/", orderRouter.ListOrders)
	// Create an order
	orderGroup.POST("/", orderRouter.CreateOrder)
	// Get an order
	orderGroup.GET("/{id}", orderRouter.LoadOrder)
	// Confirm an order. This also renew or upgrade
	// membership.
	orderGroup.PATCH("/{id}", orderRouter.ConfirmOrder)

	promoRouter := controller.NewPromoRouter(db)
	promoGroup := apiBase.Group("/promos")
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
	androidGroup := apiBase.Group("/android")
	androidGroup.GET("/exists/:versionName", androidRouter.TagExists)
	androidGroup.POST("/releases", androidRouter.CreateRelease)
	androidGroup.GET("/releases", androidRouter.Releases)
	androidGroup.GET("/releases/:versionName", androidRouter.SingleRelease)
	androidGroup.PATCH("/releases/:versionName", androidRouter.UpdateRelease)
	androidGroup.DELETE("/releases/:versionName", androidRouter.DeleteRelease)

	statsRouter := controller.NewStatsRouter(db)
	statsGroup := apiBase.Group("/stats")
	statsGroup.GET("/signup/daily", statsRouter.DailySignUp)
	statsGroup.GET("/income/year/{year}", statsRouter.YearlyIncome)

	// Search
	searchGroup := apiBase.Group("/search")
	// Search by cms user's name: /search/staff?name=<user_name>
	searchGroup.GET("/staff", staffRouter.Search)
	// Search ftc account: /search/reader?q=<email>&kind=ftc
	// Search wx account: /search/reader?q=<nickname>&kind=<wechat>&page=<number>&per_page=<number>
	searchGroup.GET("/reader", readerRouter.SearchAccount)

	//apnRouter := controller.NewAPNRouter(apnDB)
	//contentRouter := controller.NewContentRouter(db)
	//mux.Route("/apn", func(r chi.Router) {
	//
	//	r.Route("/latest", func(r chi.Router) {
	//		r.Get("/story", contentRouter.LatestStoryList)
	//	})
	//
	//	r.Route("/search", func(r chi.Router) {
	//		r.Get("/story/{id}", contentRouter.StoryTeaser)
	//		r.Get("/video/{id}", contentRouter.VideoTeaser)
	//		r.Get("/gallery/{id}", contentRouter.GalleryTeaser)
	//		r.Get("/interactive/{id}", contentRouter.InteractiveTeaser)
	//	})
	//
	//	r.Route("/stats", func(r chi.Router) {
	//		r.Get("/messages", apnRouter.ListMessages)
	//		r.Get("/timezones", apnRouter.LoadTimezones)
	//		r.Get("/devices", apnRouter.LoadDeviceDist)
	//		r.Get("/invalid", apnRouter.LoadInvalidDist)
	//	})
	//
	//	r.Route("/test-devices", func(r chi.Router) {
	//		r.Get("/", apnRouter.ListTestDevice)
	//		r.Post("/", apnRouter.CreateTestDevice)
	//		r.Delete("/{id}", apnRouter.RemoveTestDevice)
	//	})
	//})

	e.Logger.Fatal(e.Start(":3100"))

	//logger.Info("Server starts on port 3100")
	//log.Fatal(http.ListenAndServe(":3100", mux))
}
