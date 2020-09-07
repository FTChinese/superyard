package controller

import (
	"strconv"
)

func ParseInt(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 0)
}
