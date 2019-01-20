package controller

import (
	"net/http"
	"strings"

	"github.com/FTChinese/go-rest/view"
	log "github.com/sirupsen/logrus"
)

const userNameKey = "X-User-Name"

// CheckUserName middleware makes sure all request header contains `X-User-Name` field.
//
// - 401 Unauthorized if request header does not have `X-User-Name`,
// or the value is empty.
func CheckUserName(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userName := req.Header.Get(userNameKey)

		userName = strings.TrimSpace(userName)
		if userName == "" {
			log.WithField("location", "middleware: checkUserName").Info("Missing X-User-Name header")

			view.Render(w, view.NewUnauthorized(""))

			return
		}

		req.Header.Set(userNameKey, userName)

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

// ParseForm perform req.ParseForm and stops if any parse failed
// func ParseForm(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, req *http.Request) {
// 		err := req.ParseForm()

// 		if err != nil {
// 			view.Render(w, util.NewBadRequest(err.Error()))

// 			return
// 		}

// 		next.ServeHTTP(w, req)
// 	}

// 	return http.HandlerFunc(fn)
// }
