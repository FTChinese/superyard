package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-mail/mail"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/controller"
	"gitlab.com/ftchinese/backyard-api/util"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	err := godotenv.Load()
	if err != nil {
		log.WithField("package", "backyard-api.main").Error(err)

		os.Exit(1)
	}
}
func main() {

	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")

	mailHost := os.Getenv("MAILER_HOST")
	mailUser := os.Getenv("MAILER_USER")
	mailPass := os.Getenv("MAILER_PASS")

	db, err := util.NewDB(host, port, user, pass)
	if err != nil {
		os.Exit(1)
	}

	dialer := mail.NewDialer(mailHost, 587, mailUser, mailPass)

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	staffRouter := controller.NewStaffRouter(db, dialer)

	adminRouter := controller.NewAdminRouter(db, dialer)

	ftcAPIRouter := controller.NewFTCAPIRouter(db)

	ftcUserRouter := controller.NewFTCUserRouter(db)

	statsRouter := controller.NewStatsRouter(db)

	// staff router performs user login related tasks
	mux.Route("/staff", func(r1 chi.Router) {

		r1.Post("/auth", staffRouter.Auth)

		r1.Route("/password-reset", func(r2 chi.Router) {
			r2.Post("/", staffRouter.ResetPassword)

			r2.Post("/letter", staffRouter.ForgotPassword)

			r2.Get("/tokens/{token}", staffRouter.VerifyToken)
		})
	})

	mux.Route("/user", func(r chi.Router) {
		r.Use(controller.CheckUserName)

		r.Get("/profile", staffRouter.Profile)

		r.Patch("/display-name", staffRouter.UpdateDisplayName)

		r.Patch("/email", staffRouter.UpdateEmail)

		r.Patch("/password", staffRouter.UpdatePassword)

		r.Get("/myft", staffRouter.ListMyft)

		r.Post("/myft", staffRouter.AddMyft)

		r.Delete("/myft/{id}", staffRouter.DeleteMyft)
	})

	mux.Route("/admin", func(r chi.Router) {
		r.Use(controller.CheckUserName)

		r.Route("/staff", func(r2 chi.Router) {
			r2.Get("/exists", adminRouter.Exists)

			r2.Post("/new", adminRouter.NewStaff)

			r2.Get("/roster", adminRouter.StaffRoster)

			r2.Get("/profile/{name}", adminRouter.StaffProfile)

			r2.Put("/profile/{name}", adminRouter.ReinstateStaff)

			r2.Patch("/profile/{name}", adminRouter.UpdateStaff)

			r2.Delete("/profile/{name}", adminRouter.DeleteStaff)
		})

		r.Route("/vip", func(r2 chi.Router) {
			r2.Get("/", adminRouter.VIPRoster)

			r2.Put("/{myftId}", adminRouter.GrantVIP)

			r2.Delete("/{myftId}", adminRouter.RevokeVIP)
		})
	})

	mux.Route("/ftc-api", func(r chi.Router) {

		r.Use(controller.CheckUserName)

		r.Route("/apps", func(r2 chi.Router) {
			r2.Post("/", ftcAPIRouter.NewApp)

			r2.Get("/", ftcAPIRouter.ListApps)

			r2.Get("/{name}", ftcAPIRouter.GetApp)

			r2.Patch("/{name}", ftcAPIRouter.UpdateApp)

			r2.Delete("/{name}", ftcAPIRouter.DeleteApp)

			r2.Post("/{name}/transfer", ftcAPIRouter.TransferApp)
		})

		r.Route("/tokens", func(r2 chi.Router) {
			r2.Post("/", ftcAPIRouter.NewToken)

			r2.Get("/personal", ftcAPIRouter.PersonalTokens)

			r2.Delete("/personal/{tokenId}", ftcAPIRouter.DeletePersonalToken)

			r2.Get("/app/{name}", ftcAPIRouter.AppTokens)

			r2.Delete("/app/{name}/{tokenId}", ftcAPIRouter.DeleteAppToken)
		})
	})

	mux.Route("/search", func(r chi.Router) {
		r.Get("/user", ftcUserRouter.SearchUser)
	})

	mux.Route("/ftc-user", func(r chi.Router) {
		r.Route("/profile", func(r2 chi.Router) {
			r2.Get("/{userId}", ftcUserRouter.UserProfile)
			r2.Get("/{userId}/orders", ftcUserRouter.UserOrders)
			r2.Get("/{userId}/login", ftcUserRouter.LoginHistory)
		})
	})

	mux.Route("/stats", func(r chi.Router) {
		r.Get("/signup/daily", statsRouter.DailySignup)
	})

	log.Fatal(http.ListenAndServe(":3100", mux))
}
