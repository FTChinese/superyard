package controller

import (
	"errors"
	"net/http"
	"strconv"
)

type queryValue string

func (v queryValue) isEmpty() bool {
	return string(v) == ""
}

func (v queryValue) toInt() (uint, error) {
	if v.isEmpty() {
		return 0, errors.New("query: empty value")
	}

	num, err := strconv.ParseUint(string(v), 10, 0)

	if err != nil {
		return 0, err
	}

	return uint(num), nil
}

func (v queryValue) toBool() (bool, error) {
	if v.isEmpty() {
		return false, errors.New("query: empty value")
	}

	b, err := strconv.ParseBool(string(v))

	if err != nil {
		return false, err
	}

	return b, nil
}

func getQuery(req *http.Request, key string) queryValue {
	value := req.Form.Get(key)

	return queryValue(value)
}
