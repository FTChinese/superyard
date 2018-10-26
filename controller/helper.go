package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/ftchinese/backyard-api/util"
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

// Convert paramValue to boolean value.
// Returns error if the paramValue cannot be converted.
func (v paramValue) toBool() (bool, error) {
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

// Normalize the format and order of start and end time.
// Format is YYYY-MM-DD.
// Start must be before end and will be reversed if not.
// If both of them are empty, start will default to 7 days ago and end to now.
func normalizeTimeRange(start, end string) (string, string, error) {
	// If neither start nor end is supplied.
	if start == "" && end == "" {
		now := time.Now()
		start = util.SQLDateFormatter.FromTime(now.AddDate(0, 0, -7))
		end = util.SQLDateFormatter.FromTime(now)

		return start, end, nil
	}

	// If only end supplied
	if start == "" {
		endTime, err := util.ParseSQLDate(end)
		if err != nil {
			return "", "", err
		}

		startTime := endTime.AddDate(0, 0, -7)

		start = util.SQLDateFormatter.FromTime(startTime)

		return start, end, nil
	}

	if end == "" {
		startTime, err := util.ParseSQLDate(start)
		if err != nil {
			return "", "", err
		}

		endTime := startTime.AddDate(0, 0, 7)

		end = util.SQLDateFormatter.FromTime(endTime)

		return start, end, nil
	}

	startTime, err := util.ParseSQLDate(start)
	endTime, err := util.ParseSQLDate(end)

	if err != nil {
		return "", "", err
	}

	if startTime.After(endTime) {
		return end, start, nil
	}

	return start, end, nil
}
