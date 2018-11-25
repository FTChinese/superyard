package subscription

import (
	"encoding/json"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Schedule represents the beginning and ending time of
// a promotion event.
type Schedule struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Start       string `json:"startAt"`
	End         string `json:"endAt"`
	CreatedAt   string `json:"createdAt"`
	CreatedBy   string `json:"createdBy"`
}

// NewSchedule saves a new promotion schedule.
// Return the inserted row's id so that client knows which row to update in the following step.
func (env Env) NewSchedule(s Schedule) (int64, error) {
	query := `
	INSERT INTO premium.promotion_schedule
	SET name = ?,
		description = ?,
		start_utc = ?,
		end_utc = ?,
		created_utc = UTC_TIMESTAMP(),
		created_by = ?`

	start := util.SQLDatetimeUTC.FromISO8601(s.Start)
	end := util.SQLDatetimeUTC.FromISO8601(s.End)

	result, err := env.DB.Exec(query,
		s.Name,
		s.Description,
		start,
		end,
		s.CreatedBy,
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

// NewPromo saves a new discount schedule into database.
func (env Env) NewPromo(s Promotion) error {
	query := `
	INSERT INTO premium.discount_schedule
	SET name = ?,
		description = ?,
		start_utc = ?,
		end_utc = ?,
		plans = ?,
		created_utc = UTC_TIMESTAMP(),
		created_by = ?`

	startUTC := util.SQLDatetimeUTC.FromISO8601(s.Start)
	endUTC := util.SQLDatetimeUTC.FromISO8601(s.End)
	plans, err := json.Marshal(s.Plans)

	if err != nil {
		logger.WithField("location", "NewSchedule").Error(err)
		return err
	}

	_, err = env.DB.Exec(query,
		s.Name,
		s.Description,
		startUTC,
		endUTC,
		string(plans),
		s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil
}
