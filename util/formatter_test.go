package util

import (
	"testing"
	"time"
)

func TestDatetime(t *testing.T) {
	parsedTime, err := time.Parse(ISO9075, "2012-08-24 14:32:38")

	if err != nil {
		t.Error(err)
	}

	result := parsedTime.UTC().Format(time.RFC3339)

	t.Log(result)
}

func TestISO8601(t *testing.T) {
	r := ISO8601UTC.FromNow()
	t.Logf("ISO8601 now: %s", r)

	r = ISO8601UTC.FromNowDays(7)
	t.Logf("ISO8601 7 days later: %s", r)

	r = ISO8601UTC.FromDatetime("2012-08-24 14:32:38", TZShanghai)
	t.Logf("ISO8601 from date time: %s", r)

	r = SQLDateUTC8.FromNow()
	t.Logf("SQL date for now in utc8: %s", r)

	r = SQLDateUTC8.FromNowDays(7)
	t.Logf("SQL date for 7 days later in utc8: %s", r)
}
