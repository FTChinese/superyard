package controller

import (
	"github.com/gorilla/schema"
	"strconv"
)

func ParseInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 0)
}

var decoder = schema.NewDecoder()
