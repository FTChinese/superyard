package main

import (
	"flag"
	"fmt"
	"github.com/FTChinese/go-rest/view"
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
	version string
	build   string
	logger  = log.WithField("project", "backyard-api").WithField("package", "main")
	config  = util.BuildConfig{
		Version: version,
		BuiltAt: build,
	}
)

func init() {
	flag.BoolVar(&config.IsProduction, "production", false, "Indicate productions environment if present")
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
	if config.IsProduction {
		err = viper.UnmarshalKey("mysql.master", &dbConn)
	} else {
		err = viper.UnmarshalKey("mysql.dev", &dbConn)
	}

	if err != nil {
		logger.WithField("trace", "main").Error(err)
		os.Exit(1)
	}

	if config.IsProduction {
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

	searchRouter := controller.NewSearchRouter(db)
	apiRouter := controller.APIRouter(db)

	readerRouter := controller.NewReaderRouter(db)
	memberRouter := controller.NewMemberRouter(db)
	orderRouter := controller.NewOrderRouter(db)
	promoRouter := controller.NewPromoRouter(db)

	statsRouter := controller.NewStatsRouter(db)

	apnRouter := controller.NewAPNRouter(apnDB)
	contentRouter := controller.NewContentRouter(db)
	androidRouter := controller.NewAndroidRouter(db)

	// Input {userName: string, password: string}
	mux.Post("/login", staffRouter.Login)
	mux.Route("/password-reset", func(r chi.Router) {
		r.Post("/", staffRouter.ResetPassword)

		r.Post("/letter", staffRouter.ForgotPassword)

		r.Get("/tokens/{token}", staffRouter.VerifyToken)
	})

	mux.Route("/search", func(r chi.Router) {
		// /staff?email=<name@ftchinese.com>
		// /staff?name=<user_name>
		r.Get("/staff", searchRouter.Staff)
		// /reader/ftc?email=<email@example.org>
		r.Get("/reader/ftc", searchRouter.SearchFtcUser)
		// /reader/wx?q=<nickname>&page=<int>&per_page=<int>
		r.Get("/reader/wx", searchRouter.SearchWxUser)
	})

	mux.Route("/staff", func(r chi.Router) {

		//	GET /staff?page=<number>&per_page=<number>
		r.Get("/", staffRouter.List)

		// Create a staff
		// 	POST /staff
		r.Post("/", staffRouter.Create)

		// Get the staff profile
		r.Get("/{id}", staffRouter.Profile)

		// UpdateProfile a staff's profile
		r.Patch("/{id}", staffRouter.Update)

		// Delete a staff.
		r.Delete("/{id}", staffRouter.Delete)

		// Reinstate a deactivated staff
		r.Put("/{id}", staffRouter.Reinstate)

		r.Patch("/{id}/password", staffRouter.UpdatePassword)

		// NOTE: the following way of router does not work.
		// It makes GET /staff/id not responding.
		//r.Route("/{id}", func(r chi.Router) {
		//
		//	r.Patch("/display-name", staffRouter.UpdateDisplayName)
		//
		//	r.Patch("/email", staffRouter.UpdateEmail)
		//
		//	r.Patch("/password", staffRouter.UpdatePassword)
		//
		//	r.Post("/reinstate", staffRouter.Reinstate)
		//})
	})

	// Handle VIPs
	mux.Route("/vip", func(r chi.Router) {
		r.Get("/", readerRouter.ListVIP)

		r.Put("/{id}", readerRouter.GrantVIP)

		r.Delete("/{id}", readerRouter.RevokeVIP)
	})

	mux.Route("/readers", func(r chi.Router) {

		r.Route("/ftc", func(r chi.Router) {
			// FTC account
			r.Get("/{id}", readerRouter.LoadFTCAccount)
			r.Get("/{id}/profile", readerRouter.LoadFtcProfile)
			// Login history
			r.Get("/{id}/login", readerRouter.LoadLoginHistory)
		})

		r.Route("/wx", func(r chi.Router) {
			// Wx Account
			r.Get("/{id}", readerRouter.LoadWxAccount)
			r.Get("/{id}/profile", readerRouter.LoadWxProfile)
			// Wx login history
			r.Get("/{id}/login", readerRouter.LoadOAuthHistory)
		})
	})

	mux.Route("/orders", func(r chi.Router) {
		// Get a list of orders of a specific reader.
		// /orders?ftc_id=<string>&union_id=<string>&page=<int>&per_page=<int>
		// ftc_id and union_id are not both required,
		// but at least one should be present.
		r.Get("/", orderRouter.ListOrders)
		// Create an order
		r.Post("/", orderRouter.CreateOrder)
		// Get an order
		r.Get("/{id}", orderRouter.LoadOrder)
		// Confirm an order. This also renew or upgrade
		// membership.
		r.Patch("/{id}", orderRouter.ConfirmOrder)
	})

	mux.Route("/memberships", func(r chi.Router) {
		// Ge a list of memberships
		r.Get("/", memberRouter.ListMembers)
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
		r.Post("/", memberRouter.CreateMember)
		// Get one subscription
		r.Get("/{id}", memberRouter.LoadMember)
		// UpdateProfile a subscription
		r.Patch("/{id}", memberRouter.UpdateMember)
		// Delete a subscription
		r.Delete("/{id}", memberRouter.DeleteMember)
	})

	mux.Route("/promos", func(r chi.Router) {
		// ListStaff promos by page
		r.Get("/", promoRouter.ListPromos)

		// Create a new promo
		r.Post("/", promoRouter.CreateSchedule)

		// Get a promo
		r.Get("/{id}", promoRouter.LoadPromo)

		// Delete a promo
		r.Delete("/{id}", promoRouter.DisablePromo)

		r.Patch("/{id}/plans", promoRouter.SetPricingPlans)
		r.Patch("/{id}/banner", promoRouter.SetBanner)
	})

	mux.Route("/api", func(r chi.Router) {

		r.Route("/apps", func(r chi.Router) {
			r.Get("/", apiRouter.ListApps)

			r.Post("/", apiRouter.CreateApp)

			r.Get("/{id}", apiRouter.LoadApp)

			r.Patch("/{id}", apiRouter.UpdateApp)

			r.Delete("/{id}", apiRouter.RemoveApp)
		})

		r.Route("/keys", func(r chi.Router) {
			// /api/keys?client_id=<string>&page=<number>&per_page=<number>
			// /api/keys?staff_name=<string>&page=<number>&per_page=<number>
			r.Get("/", apiRouter.ListKeys)

			// Create a new key.
			r.Post("/", apiRouter.CreateKey)

			// Delete all keys owned by someone.
			// You cannot delete all keys belonging to an app
			// here since it is performed when an app is deleted.
			r.Delete("/", apiRouter.DeletePersonalKeys)

			// Delete a single key belong to an app or a human
			r.Delete("/{id}", apiRouter.RemoveKey)
		})
	})

	mux.Route("/stats", func(r chi.Router) {

		r.Get("/signup/daily", statsRouter.DailySignUp)

		r.Get("/income/year/{year}", statsRouter.YearlyIncome)
	})

	mux.Route("/android", func(r chi.Router) {

		r.Get("/exists/{versionName}", androidRouter.TagExists)
		r.Post("/releases", androidRouter.CreateRelease)
		r.Get("/releases", androidRouter.Releases)
		r.Get("/releases/{versionName}", androidRouter.SingleRelease)
		r.Patch("/releases/{versionName}", androidRouter.UpdateRelease)
		r.Delete("/releases/{versionName}", androidRouter.DeleteRelease)
	})

	mux.Route("/apn", func(r chi.Router) {

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

	mux.Get("/__version", func(writer http.ResponseWriter, request *http.Request) {
		_ = view.Render(writer, view.NewResponse().NoCache().SetBody(config))
	})

	logger.Info("Server starts on port 3100")
	log.Fatal(http.ListenAndServe(":3100", mux))
}
