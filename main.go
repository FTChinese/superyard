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
	"gitlab.com/ftchinese/backyard-api/util"
)

var (
	isProd  bool
	version string
	build   string
	logger  = log.WithField("project", "backyard-api").WithField("package", "main")
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

	// Get email server config.
	var emailConn util.Conn
	err = viper.UnmarshalKey("email.ftc", &emailConn)
	if err != nil {
		logger.WithField("trace", "main").Error(err)
		os.Exit(1)
	}

	db, err := util.NewDB(dbConn)
	if err != nil {
		log.WithField("package", "backyard-api.main").Error(err)
		os.Exit(1)
	}
	logger.
		WithField("trace", "main").
		Infof("Connected to MySQL server %s", dbConn.Host)

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

	mux.Get("/__version", controller.Version(version, build))

	mux.Post("/login", staffRouter.Login)
	mux.Route("/password-reset", func(r chi.Router) {
		r.Post("/", staffRouter.ResetPassword)

		r.Post("/letter", staffRouter.ForgotPassword)

		r.Get("/tokens/{token}", staffRouter.VerifyToken)
	})

	mux.Route("/staff", func(r chi.Router) {
		r.Use(controller.StaffName)

		r.Get("/account", staffRouter.Account)

		r.Get("/profile", staffRouter.Profile)

		r.Patch("/display-name", staffRouter.UpdateDisplayName)

		r.Patch("/email", staffRouter.UpdateEmail)

		r.Patch("/password", staffRouter.UpdatePassword)

		r.Post("/myft", staffRouter.AddMyft)

		r.Get("/myft", staffRouter.ListMyft)

		r.Delete("/myft", staffRouter.DeleteMyft)
	})

	mux.Route("/admin", func(r chi.Router) {
		// TODO: add `X-Admin-Name` for access control.
		r.Use(controller.StaffName)

		r.Route("/account", func(r chi.Router) {
			r.Get("/exists", adminRouter.Exists)
			r.Get("/search", adminRouter.FindAccount)
		})

		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", adminRouter.CreateAccount)

			// GET /admin/accounts?page=<number>&per_page=<number>
			r.Get("/", adminRouter.ListAccounts)

			r.Get("/{name}", adminRouter.StaffProfile)

			r.Patch("/{name}", adminRouter.UpdateAccount)

			r.Delete("/{name}", adminRouter.DeleteStaff)

			r.Put("/{name}", adminRouter.ReinstateStaff)
		})

		r.Route("/vip", func(r2 chi.Router) {
			r2.Get("/", adminRouter.ListVIP)

			r2.Put("/{email}", adminRouter.GrantVIP)

			r2.Delete("/{email}", adminRouter.RevokeVIP)
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

			r.Post("/{name}/transfer", nextAPIRouter.TransferApp)

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
	})

	logger.Info("Server starts on port 3100")
	log.Fatal(http.ListenAndServe(":3100", mux))
}
