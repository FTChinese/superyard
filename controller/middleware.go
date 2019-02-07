package controller

import (
	"net/http"
	"strings"

	"github.com/FTChinese/go-rest/view"
	log "github.com/sirupsen/logrus"
)

const (
	staffNameKey = "X-Staff-Name"
	adminNameKey = "X-Admin-Name"
	userIDKey    = "X-User-Id"
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
		userName := req.Header.Get(staffNameKey)

		userName = strings.TrimSpace(userName)
		if userName == "" {
			log.WithField("trace", "middleware: checkUserName").Info("Missing X-User-Name header")

			view.Render(w, view.NewUnauthorized(""))

			return
		}

		req.Header.Set(staffNameKey, userName)

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

func UserID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userID := req.Header.Get(userIDKey)

		userID = strings.TrimSpace(userID)
		if userID == "" {
			log.WithField("trace", "UserID").Info("Missing X-User-Name header")

			view.Render(w, view.NewUnauthorized(""))

			return
		}

		req.Header.Set(userIDKey, userID)

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}