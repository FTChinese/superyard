package util

import (
	"log"
	"time"
)

// MySQL time layout
const (
	ISO9075     = "2006-01-02 15:04:05" // Layout for SQL DATETIME
	ISO9075Date = "2006-01-02"          // Layout for SQL DATE
)

const (
	secondsOfMinute = 60
	secondsOfHour   = 60 * secondsOfMinute
)

// Fixed time zone
var (
	TZShanghai = time.FixedZone("UTC+8", 8*secondsOfHour)
)

// Formatter instances
var (
	// ISO8601UTC formats time to RFC3339 in UTC
	ISO8601UTC = Formatter{time.RFC3339, time.UTC}
	// SQLDatetimeUTC formats time to SQL DATETIME in UTC
	SQLDatetimeUTC = Formatter{ISO9075, time.UTC}
	// SQLDateForamtter formats time to SQL DATE in UTC+08
	SQLDateUTC8 = Formatter{ISO9075Date, TZShanghai}
)

// Formatter converts a time.Time instance to specified layout in specified location
type Formatter struct {
	layout string         // target layout
	loc    *time.Location // target timezone
}

// ToLocation returns a new Formatter with location changed to the specified time zone
func (f Formatter) ToLocation(loc *time.Location) Formatter {
	f.loc = loc
	return f
}

// FromUnix formats a Unix timestamp to human readable string
func (f Formatter) FromUnix(sec int64) string {
	return time.Unix(sec, 0).In(f.loc).Format(f.layout)
}

// FromDatetime parsed SQL DATETIME and converts to specified format
func (f Formatter) FromDatetime(value string, loc *time.Location) string {
	if loc == nil {
		loc = time.UTC
	}

	t, err := time.ParseInLocation(ISO9075, value, loc)

	log.Println(t)

	if err != nil {
		return value
	}
	return t.In(f.loc).Format(f.layout)
}

func (f Formatter) FromISO8601(value string) string {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}

	return t.In(f.loc).Format(f.layout)
}

// FromNow converts current time to specified format
func (f Formatter) FromNow() string {
	return time.Now().In(f.loc).Format(f.layout)
}

// FromNowDays converts a future or past date to the specified format
func (f Formatter) FromNowDays(days int) string {
	return time.Now().AddDate(0, 0, days).In(f.loc).Format(f.layout)
}

// FromTime formats a time.Time instance to the specified layout.
func (f Formatter) FromTime(t time.Time) string {
	return t.In(f.loc).Format(f.layout)
}

// ParseSQLDate return time.Time by parsing a MySQL DATE value.
func ParseSQLDate(value string) (time.Time, error) {
	return time.Parse(ISO9075Date, value)
}
