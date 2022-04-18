package controller

import (
	"github.com/gorilla/schema"
	"net/http"
)

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

type LiveRefresh struct {
	Live    bool `schema:"live"`
	Refresh bool `schema:"refresh"`
}
