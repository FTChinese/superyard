package subscription

import (
	"strings"

	"github.com/FTChinese/go-rest/view"
	"gitlab.com/ftchinese/backyard-api/util"
)

// Schedule represents the beginning and ending time of
// a promotion event.
type Schedule struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`        // Required. Max 256 chars
	Description string `json:"description"` // Optional. Max 256 chars
	Start       string `json:"startAt"`     // Required. ISO 8601 date time string.
	End         string `json:"endAt"`       // Required. ISO 8601 date time string.
}

// Sanitize removes leading and trailing spaces of each string fields.
func (s *Schedule) Sanitize() {
	s.Name = strings.TrimSpace(s.Name)
	s.Description = strings.TrimSpace(s.Description)
	s.Start = strings.TrimSpace(s.Start)
	s.End = strings.TrimSpace(s.End)
}

// Validate validates incoming data for a new schedule.
func (s *Schedule) Validate() *view.Reason {
	if r := util.RequireNotEmptyWithMax(s.Name, 256, "name"); r != nil {
		return r
	}

	if r := util.OptionalMaxLen(s.Description, 256, "description"); r != nil {
		return r
	}

	if r := util.RequireNotEmpty(s.Start, "startAt"); r != nil {
		return r
	}

	return util.RequireNotEmpty(s.End, "endAt")
}

// NewSchedule saves a new promotion schedule.
// Return the inserted row's id so that client knows which row to update in the following step.
func (env Env) NewSchedule(s Schedule, creator string) (int64, error) {
	query := `
	INSERT INTO premium.promotion_schedule
	SET name = ?,
		description = ?,
		start_utc = ?,
		end_utc = ?,
		created_by = ?`

	start := util.SQLDatetimeUTC.FromISO8601(s.Start)
	end := util.SQLDatetimeUTC.FromISO8601(s.End)

	result, err := env.DB.Exec(query,
		s.Name,
		s.Description,
		start,
		end,
		creator,
	)

	if err != nil {
		logger.WithField("location", "NewSchedule").Error(err)
		return -1, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		logger.WithField("location", "NewSchedule").Error(err)
		return -1, err
	}

	return id, nil
}
