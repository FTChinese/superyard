package util

import "time"

const (
	iso8601 = "2006-01-02T15:04:05Z" // Layout for ISO8601
	iso9075 = "2006-01-02 15:04:05"  // Layout for MySQL DATETIME
)

var utc8 = time.FixedZone("UTC+8", 8*60*60)

// FormatUTCDatetime converts mysql DATETIME in utc zone to ISO8601 in UTC
func FormatUTCDatetime(value string) string {
	t, err := time.Parse(iso9075, value)

	if err != nil {
		return value
	}

	return t.UTC().Format(time.RFC3339)
}

// FormatUTC8Datetime converts mysql DATETIME in +08:00 offset to ISO8601 in UTC
func FormatUTC8Datetime(value string) string {
	t, err := time.ParseInLocation(iso9075, value, utc8)

	if err != nil {
		return value
	}

	return t.UTC().Format(time.RFC3339)
}

// FormatUnix turns unix timestamp in seconds to ISO8601 in UTC.
func FormatUnix(sec int64) string {
	return time.Unix(sec, 0).UTC().Format(time.RFC3339)
}
