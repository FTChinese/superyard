package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var logger = logrus.
	WithField("project", "superyard").
	WithField("package", "controller")

// GetURLParam gets a url parameter.
func GetURLParam(req *http.Request, key string) gorest.Param {
	v := chi.URLParam(req, key)

	return gorest.NewParam(key, v)
}

func ParseInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 0)
}
