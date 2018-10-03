package controller

import (
	"errors"
	"net/http"
	"strconv"
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
	if v.isEmpty() {
		return false, errors.New("query: empty value")
	}

	b, err := strconv.ParseBool(string(v))

	if err != nil {
		return false, err
	}

	return b, nil
}

func getQueryParam(req *http.Request, key string) paramValue {
	value := req.Form.Get(key)

	return paramValue(value)
}
