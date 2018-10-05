package util

import "time"

const (
	iso8601     = "2006-01-02T15:04:05Z" // Layout for ISO8601
	iso9075     = "2006-01-02 15:04:05"  // Layout for MySQL DATETIME
	iso9075Date = "2006-01-02"
)

var utc8Offset = 8 * 60 * 60
var utc8 = time.FixedZone("UTC+8", utc8Offset)
var secondsOfDay = 24 * 60 * 60

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

// DateNow gets a string represetation of current date in YYYY-MM-DD format.
// By specifying a duration to add you can get it at any time zone
// Example:
// To get current time in +08:00, DateNow(8*60*60)
// To get 7 days from now in +08:00, DateNow(8*60*60 + 7*24*60*60)
func DateNow(offset int) string {
	return time.Now().UTC().Add(time.Duration(offset)).Format(iso9075Date)
}

type Time struct {
	time.Time
}

func Now() Time {
	t := Time{}
	t.UTC()

	return t
}

func ParseISO9075(value string, zone int) (Time, error) {
	if zone == 0 {
		t, err := time.Parse(iso9075, value)
		if err != nil {
			return Time{}, err
		}
		return Time{t.UTC()}, nil
	}

	t, err := time.ParseInLocation(iso9075, value, time.FixedZone("UTC+8", zone*60*60))

	if err != nil {
		return Time{}, err
	}
	return Time{t.UTC()}, nil
}

func Unix(sec int64) Time {
	return Time{time.Unix(sec, 0).UTC()}
}

func (t Time) plus(offset int) Time {
	t.Add(time.Duration(offset))
	return t
}

func (t Time) PlusYears(years int) Time {
	t.plus(t.YearDay() * secondsOfDay)
	return t
}

func (t Time) PlusDays(days int) Time {
	t.plus(days * secondsOfDay)
	return t
}

func (t Time) format(layout string, zone int) string {
	return t.Add(time.Duration(zone * secondsOfDay)).Format(layout)
}

func (t Time) FormatISO8601(zone int) string {
	return t.format(time.RFC3339, zone)
}

func (t Time) FormatISO9075(zone int) string {
	return t.format(iso9075, zone)
}

func (t Time) FormatISO9075Date(zone int) string {
	return t.format(iso9075Date, zone)
}
