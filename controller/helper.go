package controller

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gitlab.com/ftchinese/backyard-api/types/util"
	"io"
	"io/ioutil"
	"net/http"
)

var logger = log.WithField("project", "backyard-api").WithField("package", "controller")

// GetURLParam gets a url parameter.
func GetURLParam(req *http.Request, key string) gorest.Param {
	v := chi.URLParam(req, key)

	return gorest.NewParam(key, v)
}

func GetJSONResult(data io.ReadCloser, path string) (gjson.Result, error) {
	b, err := ioutil.ReadAll(data)
	defer data.Close()

	if err != nil {
		return gjson.Result{}, err
	}

	return gjson.GetBytes(b, path), nil
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
