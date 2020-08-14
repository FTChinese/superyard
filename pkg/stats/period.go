package stats

import (
	"time"

	"github.com/FTChinese/go-rest/chrono"
)

// Period represents a range of time.
type Period struct {
	Start chrono.Time
	End   chrono.Time
}

// NewPeriod creates a Period from two date strings.
// Returns error if the data string cannot be parsed.
// Format: YYYY-MM-DD
// Start should be before end, but the order will be corrected if reversed.
func NewPeriod(start, end string) (Period, error) {
	// If neither start nor end is supplied.
	if start == "" && end == "" {
		now := time.Now()

		return Period{
			Start: chrono.TimeFrom(now.AddDate(0, 0, -7)),
			End:   chrono.TimeFrom(now),
		}, nil
	}

	// If only end supplied
	if start == "" {
		endTime, err := chrono.ParseDateTime(end, chrono.TZShanghai)
		if err != nil {
			return Period{}, err
		}

		return Period{
			Start: chrono.TimeFrom(endTime.AddDate(0, 0, -7)),
			End:   chrono.TimeFrom(endTime),
		}, nil
	}

	if end == "" {
		startTime, err := chrono.ParseDateTime(start, chrono.TZShanghai)
		if err != nil {
			return Period{}, err
		}

		return Period{
			Start: chrono.TimeFrom(startTime),
			End:   chrono.TimeFrom(startTime.AddDate(0, 0, 7)),
		}, nil
	}

	startTime, err := chrono.ParseDateTime(start, chrono.TZShanghai)
	endTime, err := chrono.ParseDateTime(end, chrono.TZShanghai)

	if err != nil {
		return Period{}, err
	}

	if startTime.After(endTime) {
		return Period{
			Start: chrono.TimeFrom(endTime),
			End:   chrono.TimeFrom(startTime),
		}, nil
	}

	return Period{
		Start: chrono.TimeFrom(startTime),
		End:   chrono.TimeFrom(endTime),
	}, nil
}
