package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/staff"
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

	staffEnv := staff.Env{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/staff", func(r chi.Router) {
		r.Post("/auth", func(w http.ResponseWriter, req *http.Request) {
			var login staff.Login
			dec := json.NewDecoder(req.Body)
			err := dec.Decode(login)

			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			account, err := staffEnv.Auth(login)

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			enc := json.NewEncoder(w)
			enc.SetEscapeHTML(false)
			enc.SetIndent("", "\t")

			err = enc.Encode(account)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
	})

	log.Fatal(http.ListenAndServe(":3100", r))
}
