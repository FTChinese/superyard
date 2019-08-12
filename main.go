package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/FTChinese/go-rest/postoffice"
	"github.com/spf13/viper"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/controller"
	"gitlab.com/ftchinese/backyard-api/models/util"
)

var (
	isProd  bool
	version string
	build   string
	logger  = log.WithField("project", "backyard-api").WithField("package", "main")
	config  util.BuildConfig
)

func init() {
	flag.BoolVar(&isProd, "production", false, "Indicate productions environment if present")
	var v = flag.Bool("v", false, "print current version")

	flag.Parse()

	if *v {
		fmt.Printf("%s\nBuild at %s\n", version, build)
		os.Exit(0)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")
	err := viper.ReadInConfig()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	// Get DB connection config.
	var dbConn util.Conn
	var apnDBConn util.Conn
	var err error
	if isProd {
		err = viper.UnmarshalKey("mysql.master", &dbConn)
	} else {
		err = viper.UnmarshalKey("mysql.dev", &dbConn)
	}

	if err != nil {
		logger.WithField("trace", "main").Error(err)
		os.Exit(1)
	}

	if isProd {
		err = viper.UnmarshalKey("mysql.apn", &apnDBConn)
	} else {
		apnDBConn = dbConn
	}

	if err != nil {
		logger.WithField("trace", "main").Error(err)
		os.Exit(1)
	}

	// Get email server config.
	var emailConn util.Conn
	err = viper.UnmarshalKey("email.ftc", &emailConn)
	if err != nil {
		logger.WithField("trace", "main").Error(err)
		os.Exit(1)
	}

	db, err := util.NewDBX(dbConn)
	if err != nil {
		log.WithField("package", "backyard-api.main").Error(err)
		os.Exit(1)
	}
	logger.
		WithField("trace", "main").
		Infof("Connected to MySQL server %s", dbConn.Host)

	apnDB, err := util.NewDB(apnDBConn)
	if err != nil {
		log.WithField("package", "backyard-api.main").Error(err)
		os.Exit(1)
	}
	logger.
		WithField("trace", "main").
		Infof("Connected to MySQL APN server %s", apnDBConn.Host)

	post := postoffice.NewPostman(
		emailConn.Host,
		emailConn.Port,
		emailConn.User,
		emailConn.Pass)

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.NoCache)

	staffRouter := controller.NewStaffRouter(db, post)
	adminRouter := controller.NewAdminRouter(db, post)

	nextAPIRouter := controller.NewNextAPIRouter(db)

	userRouter := controller.NewUserRouter(db)

	statsRouter := controller.NewStatsRouter(db)

	searchRouter := controller.NewSearchRouter(db)

	subsRouter := controller.NewSubsRouter(db)

	apnRouter := controller.NewAPNRouter(apnDB)

	contentRouter := controller.NewContentRouter(db)

	androidRouter := controller.NewAndroidRouter(db)

	mux.Get("/__version", controller.Version(version, build))

	mux.Post("/login", staffRouter.Login)
	mux.Route("/password-reset", func(r chi.Router) {
		r.Post("/", staffRouter.ResetPassword)

		r.Post("/letter", staffRouter.ForgotPassword)

		r.Get("/tokens/{token}", staffRouter.VerifyToken)
	})

	mux.Route("/staff", func(r chi.Router) {
		//r.Use(controller.StaffName)

		// List all staff
		r.Get("/", staffRouter.List)

		// Create a staff
		r.Post("/", staffRouter.Create)

		// Retrieve a staff
		r.Get("/{id}", staffRouter.Profile)

		// Update a staff's profile
		r.Patch("/{id}", staffRouter.Update)

		// Delete a staff.
		r.Delete("/{id}", staffRouter.Delete)

		r.Route("/{id}", func(r chi.Router) {

			r.Patch("/display-name", staffRouter.UpdateDisplayName)
			r.Patch("/email", staffRouter.UpdateEmail)
			r.Patch("/password", staffRouter.UpdatePassword)

			r.Get("/myft", staffRouter.ListMyft)
			r.Post("/myft", staffRouter.AddMyft)
			r.Delete("/myft", staffRouter.DeleteMyft)

			r.Post("/reinstate", staffRouter.Reinstate)
		})

		//r.Get("/account", staffRouter.Account)

		//r.Get("/profile", staffRouter.Profile)

		//r.Patch("/display-name", staffRouter.UpdateDisplayName)

		//r.Patch("/email", staffRouter.UpdateEmail)

		//r.Patch("/password", staffRouter.UpdatePassword)

		//r.Post("/myft", staffRouter.AddMyft)
		//
		//r.Get("/myft", staffRouter.ListMyft)
		//
		//r.Delete("/myft", staffRouter.DeleteMyft)
	})

	mux.Route("/admin", func(r chi.Router) {

		// /admin/search/staff?k={name|email}&v={value}
		r.Route("/search", func(r chi.Router) {
			r.Get("/staff", adminRouter.SearchStaff)
		})

		r.Route("/vip", func(r chi.Router) {
			r.Get("/", adminRouter.ListVIP)

			r.Put("/{id}", adminRouter.GrantVIP)

			r.Delete("/{id}", adminRouter.RevokeVIP)
		})
	})

	mux.Route("/apn", func(r chi.Router) {
		r.Use(controller.StaffName)

		r.Route("/latest", func(r chi.Router) {
			r.Get("/story", contentRouter.LatestStoryList)
		})

		r.Route("/search", func(r chi.Router) {
			r.Get("/story/{id}", contentRouter.StoryTeaser)
			r.Get("/video/{id}", contentRouter.VideoTeaser)
			r.Get("/gallery/{id}", contentRouter.GalleryTeaser)
			r.Get("/interactive/{id}", contentRouter.InteractiveTeaser)
		})

		r.Route("/stats", func(r chi.Router) {
			r.Get("/messages", apnRouter.ListMessages)
			r.Get("/timezones", apnRouter.LoadTimezones)
			r.Get("/devices", apnRouter.LoadDeviceDist)
			r.Get("/invalid", apnRouter.LoadInvalidDist)
		})

		r.Route("/test-devices", func(r chi.Router) {
			r.Get("/", apnRouter.ListTestDevice)
			r.Post("/", apnRouter.CreateTestDevice)
			r.Delete("/{id}", apnRouter.RemoveTestDevice)
		})
	})

	mux.Route("/next", func(r chi.Router) {

		r.Use(controller.StaffName)

		r.Route("/apps", func(r chi.Router) {
			r.Post("/", nextAPIRouter.CreateApp)

			r.Get("/", nextAPIRouter.ListApps)

			r.Get("/{name}", nextAPIRouter.LoadApp)

			r.Patch("/{name}", nextAPIRouter.UpdateApp)

			r.Delete("/{name}", nextAPIRouter.RemoveApp)

			//r.Post("/{name}/transfer", nextAPIRouter.TransferApp)

			r.Post("/{name}/tokens", nextAPIRouter.NewAppToken)

			r.Get("/{name}/tokens", nextAPIRouter.ListAppTokens)

			r.Delete("/{name}/tokens/{id}", nextAPIRouter.RemoveAppToken)

		})

		r.Route("/keys", func(r chi.Router) {
			r.Post("/", nextAPIRouter.CreateKey)

			r.Get("/", nextAPIRouter.ListKeys)

			r.Delete("/{tokenId}", nextAPIRouter.RemoveKey)
		})
	})

	mux.Route("/users", func(r chi.Router) {
		r.Use(controller.StaffName)

		r.Route("/ftc", func(r chi.Router) {
			r.Get("/account/{id}", userRouter.LoadFTCAccount)
			r.Get("/orders/{id}", userRouter.LoadOrders)
			// Show login history
			r.Get("/login-history/{id}", userRouter.LoadLoginHistory)
		})

		r.Route("/wx", func(r chi.Router) {
			r.Get("/account/{id}", userRouter.LoadWxAccount)
			r.Get("/orders/{id}", userRouter.LoadOrdersWxOnly)
			r.Get("/login-history/{id}", userRouter.LoadOAuthHistory)
		})
	})

	mux.Route("/search", func(r chi.Router) {
		r.Use(controller.StaffName)

		r.Get("/user/ftc", searchRouter.SearchFTCUser)
		r.Get("/user/wx", searchRouter.SearchWxUser)
		r.Get("/order", searchRouter.SearchOrder)
		r.Get("/gift-card", searchRouter.GiftCard)
	})

	mux.Route("/subs", func(r chi.Router) {
		r.Use(controller.StaffName)

		r.Route("/promos", func(r chi.Router) {
			// List promos by page
			r.Get("/", subsRouter.ListPromos)

			// Create a new promo
			r.Post("/", subsRouter.CreateSchedule)

			// Get a promo
			r.Get("/{id}", subsRouter.LoadPromo)

			// Delete a promo
			r.Delete("/{id}", subsRouter.DisablePromo)

			r.Patch("/{id}/plans", subsRouter.SetPricingPlans)
			r.Patch("/{id}/banner", subsRouter.SetBanner)
		})
	})

	mux.Route("/stats", func(r chi.Router) {
		r.Use(controller.StaffName)

		r.Get("/signup/daily", statsRouter.DailySignUp)

		r.Get("/income/year/{year}", statsRouter.YearlyIncome)
	})

	mux.Route("/android", func(r chi.Router) {
		r.Use(controller.StaffName)

		r.Get("/exists/{versionName}", androidRouter.TagExists)
		r.Post("/releases", androidRouter.CreateRelease)
		r.Get("/releases", androidRouter.Releases)
		r.Get("/releases/{versionName}", androidRouter.SingleRelease)
		r.Patch("/releases/{versionName}", androidRouter.UpdateRelease)
		r.Delete("/releases/{versionName}", androidRouter.DeleteRelease)
	})

	logger.Info("Server starts on port 3100")
	log.Fatal(http.ListenAndServe(":3100", mux))
}
