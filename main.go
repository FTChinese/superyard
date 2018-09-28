package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"gitlab.com/ftchinese/backyard-api/staffmodel"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-sql-driver/mysql"
)

func main() {
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

	cmsUser := staffmodel.CMSUser{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/staff", func(r chi.Router) {
		r.Post("/auth", func(w http.ResponseWriter, req *http.Request) {
			var login staffmodel.StaffLogin
			dec := json.NewDecoder(req.Body)
			err := dec.Decode(login)

			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			account, err := cmsUser.Auth(login)

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
