package controller

import (
	"net/http"
	"strings"

	"github.com/FTChinese/go-rest/view"
	log "github.com/sirupsen/logrus"
)

const (
	userNameKey  = "X-User-Name"
	userEmailKey = "X-Email"
	unionIDKey   = "X-Union-Id"
)

// NoCache set Cache-Control request header
func NoCache(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Cache-Control", "must-revalidate")
		w.Header().Add("Pragma", "no-cache")
		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

// StaffName middleware makes sure all request header contains `X-User-Name` field.
//
// - 401 Unauthorized if request header does not have `X-User-Name`,
// or the value is empty.
func StaffName(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userName := req.Header.Get(userNameKey)

		userName = strings.TrimSpace(userName)
		if userName == "" {
			log.WithField("trace", "middleware: checkUserName").Info("Missing X-User-Name header")

			view.Render(w, view.NewUnauthorized(""))

			return
		}

		req.Header.Set(userNameKey, userName)

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

func FtcUserEmail(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		email := req.Header.Get(userEmailKey)

		email = strings.TrimSpace(email)
		if email == "" {
			log.WithField("trace", "FtcUserEmail").Info("Missing X-Email header")

			view.Render(w, view.NewUnauthorized(""))

			return
		}

		req.Header.Set(userEmailKey, email)

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

// Version show current version of api.
func Version(version, build string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		b := map[string]string{
			"version": version,
			"build":   build,
		}

		view.Render(w, view.NewResponse().NoCache().SetBody(b))
	}
}
