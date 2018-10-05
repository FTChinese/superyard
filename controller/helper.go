package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type paramValue string

func (v paramValue) isEmpty() bool {
	return string(v) == ""
}

func (v paramValue) toInt() (uint, error) {
	if v.isEmpty() {
		return 0, errors.New("query: empty value")
	}

	num, err := strconv.ParseUint(string(v), 10, 0)

	if err != nil {
		return 0, err
	}

	return uint(num), nil
}

func (v paramValue) toBool() (bool, error) {
	// If the paramValue does not exist, default to false value.
	if v.isEmpty() {
		return false, nil
	}

	return strconv.ParseBool(string(v))
}

func (v paramValue) toString() string {
	return string(v)
}

func getQueryParam(req *http.Request, key string) paramValue {
	value := req.Form.Get(key)

	return paramValue(value)
}

func getURLParam(req *http.Request, key string) paramValue {
	value := chi.URLParam(req, key)

	return paramValue(value)
}
