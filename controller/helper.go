package controller

import (
	"github.com/FTChinese/go-rest"
	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gitlab.com/ftchinese/backyard-api/models/util"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var logger = logrus.
	WithField("project", "backyard-api").
	WithField("package", "controller")

var decoder = schema.NewDecoder()

// GetURLParam gets a url parameter.
func GetURLParam(req *http.Request, key string) gorest.Param {
	v := chi.URLParam(req, key)

	return gorest.NewParam(key, v)
}

// GetString get a string field from http request body.
// Return empty string even if the passed in data does not contain the required key.
func GetString(data io.ReadCloser, path string) (string, error) {
	b, err := ioutil.ReadAll(data)
	defer data.Close()

	if err != nil {
		return "", err
	}

	result := gjson.GetBytes(b, path)

	if !result.Exists() {
		return "", nil
	}

	value := strings.TrimSpace(result.String())

	return value, nil
}

// IsAlreadyExists tests if an error means the field already exists
func IsAlreadyExists(err error) bool {
	if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return true
	}

	if err == util.ErrAlreadyExists {
		return true
	}

	return false
}
