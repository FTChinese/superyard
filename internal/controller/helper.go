package controller

import (
	"github.com/gorilla/schema"
	"net/http"
	"strconv"
)

func ParseInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 0)
}

var decoder = schema.NewDecoder()

func decodeForm(v interface{}, req *http.Request) error {
	decoder.IgnoreUnknownKeys(true)

	if err := req.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(v, req.Form); err != nil {
		return err
	}

	return nil
}
