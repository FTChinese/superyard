package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/controller"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	cfg := &mysql.Config{
		User:                 "sampadm",
		Passwd:               "sampadm",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		os.Exit(1)
	}

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	staffRouter := controller.NewStaffRouter(db)

	adminRouter := controller.NewAdminRouter(db)

	ftcAPIRouter := controller.NewFTCAPIRouter(db)

	// staff router performs user login related tasks
	mux.Route("/staff", func(r1 chi.Router) {
		r1.Get("/exists", staffRouter.Exists)
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
		r.Delete("/myft/{id}", staffRouter.RemoveMyft)
	})

	mux.Route("/admin", func(r chi.Router) {
		r.Use(controller.CheckUserName)

		r.Route("/staff", func(r2 chi.Router) {
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

			r2.Delete("/{name}", ftcAPIRouter.RemoveApp)

			r2.Post("/{name}/transfer", ftcAPIRouter.TransferApp)
		})

		r.Route("/tokens", func(r2 chi.Router) {
			r2.Post("/", ftcAPIRouter.NewToken)

			r2.Get("/personal", ftcAPIRouter.PersonalTokens)

			r2.Delete("/personal/{tokenId}", ftcAPIRouter.RemovePersonalToken)

			r2.Get("/app/{name}", ftcAPIRouter.AppTokens)

			r2.Delete("/app/{name}/{tokenId}", ftcAPIRouter.RemoveAppToken)
		})
	})

	log.Fatal(http.ListenAndServe(":3100", mux))
}
