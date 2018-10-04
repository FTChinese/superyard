package controller

import (
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/util"
	"gitlab.com/ftchinese/backyard-api/view"
)

const userNameKey = "X-User-Name"

// CheckUserName makes sure the request header contains `X-User-Name` fields
func CheckUserName(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userName := req.Header.Get(userNameKey)

		userName = strings.TrimSpace(userName)
		if userName == "" {
			log.WithField("location", "middleware: checkUserName").Info("Missing X-User-Name header")

			view.Render(w, util.NewUnauthorized(""))

			return
		}

		req.Header.Set(userNameKey, userName)

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

// ParseForm perform req.ParseForm and stops if any parse failed
func ParseForm(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()

		if err != nil {
			view.Render(w, util.NewBadRequest(err.Error()))

			return
		}

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
