package controller

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gitlab.com/ftchinese/backyard-api/util"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

var logger = log.WithField("project", "backyard-api").WithField("package", "controller")

// Param represents a pair of query parameter from URL.
type Param struct {
	key   string
	value string
}

// ToBool converts a query parameter to boolean value.
func (p Param) ToBool() (bool, error) {
	return strconv.ParseBool(string(p.value))
}

// ToString converts a query parameter to string value.
// Returns error for an empty value.
func (p Param) ToString() (string, error) {
	if p.value == "" {
		return "", fmt.Errorf("%s have empty value", p.key)
	}

	return p.value, nil
}

// ToInt converts the value of a query parameter to int64
func (p Param) ToInt() (int64, error) {
	if p.value == "" {
		return 0, fmt.Errorf("%s have empty value", p.key)
	}

	num, err := strconv.ParseInt(string(p.value), 10, 0)

	if err != nil {
		return 0, err
	}

	return num, nil
}

// GetQueryParam gets a pair of query parameter from URL.
func GetQueryParam(req *http.Request, key string) Param {
	v := req.Form.Get(key)

	return Param{
		key:   key,
		value: v,
	}
}

// GetURLParam gets a url parameter.
func GetURLParam(req *http.Request, key string) Param {
	v := chi.URLParam(req, key)

	return Param{
		key:   key,
		value: v,
	}
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

func GetPagination(req *http.Request) util.Pagination {
	page, _ := GetQueryParam(req, "page").ToInt()
	perPage, err := GetQueryParam(req, "per_page").ToInt()
	if err != nil {
		perPage = 20
	}

	return util.NewPagination(page, perPage)
}
