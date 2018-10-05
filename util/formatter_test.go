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
	r := ISO8601Formatter.FromNow()
	t.Logf("ISO8601 now: %s", r)

	r = ISO8601Formatter.FromNowDays(7)
	t.Logf("ISO8601 7 days later: %s", r)

	r = ISO8601Formatter.FromDatetime("2012-08-24 14:32:38", TZShanghai)
	t.Logf("ISO8601 from date time: %s", r)

	r = SQLDateFormatter.FromNow()
	t.Logf("SQL date for now in utc8: %s", r)

	r = SQLDateFormatter.FromNowDays(7)
	t.Logf("SQL date for 7 days later in utc8: %s", r)
}
