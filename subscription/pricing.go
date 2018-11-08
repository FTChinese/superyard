package subscription

import (
	"encoding/json"
	"fmt"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Plan contains details of subscription plan.
type Plan struct {
	Tier  string  `json:"tier"`
	Cycle string  `json:"cycle"`
	Price float64 `json:"price"`
	ID    int
	// For wxpay, this is used as `body` parameter;
	// For alipay, this is used as `subject` parameter.
	Description string `json:"description"` // required, max 128 chars
	// For wxpay, this is used as `detail` parameter;
	// For alipay, this is used as `body` parameter.
	Message string `json:"message"`
}

// Schedule represents a discount activity
type Schedule struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Start       string          `json:"startAt"`
	End         string          `json:"endAt"`
	Plans       map[string]Plan `json:"plans"`
	CreatedAt   string          `json:"createdAt"`
	CreatedBy   string          `json:"createdBy"`
}

// NewShedule saves a new discount schedule into database.
func (env Env) NewShedule(s Schedule) error {
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

// RetrieveSchedule loads a schedule from database.
func (env Env) RetrieveSchedule(id int64) (Schedule, error) {
	query := fmt.Sprintf(`
	%s
	WHERE id = ?
	LIMIT 1`, stmtDiscount)

	var s Schedule
	var plans string
	var start string
	var end string
	var created string
	err := env.DB.QueryRow(query, id).Scan(
		&s.ID,
		&s.Name,
		&s.Description,
		&start,
		&end,
		&plans,
		&created,
		&s.CreatedBy,
	)

	if err != nil {
		logger.WithField("location", "RetrieveSchedule").Error(err)
		return s, err
	}

	if err := json.Unmarshal([]byte(plans), &s.Plans); err != nil {
		return s, err
	}

	s.Start = util.ISO8601UTC.FromDatetime(start, nil)
	s.End = util.ISO8601UTC.FromDatetime(end, nil)
	s.CreatedAt = util.ISO8601UTC.FromDatetime(created, nil)

	return s, nil
}

// ListSchedules show all schedules.
func (env Env) ListSchedules(page, rowCount int64) ([]Schedule, error) {
	offset := (page - 1) * rowCount

	query := fmt.Sprintf(`
	%s
	ORDER BY id DESC
	LIMIT ? OFFSET ?`, stmtDiscount)

	rows, err := env.DB.Query(query, rowCount, offset)

	if err != nil {
		logger.
			WithField("location", "ListSchedules").
			Error(err)

		return nil, err
	}
	defer rows.Close()

	schs := make([]Schedule, 0)

	for rows.Next() {
		var s Schedule
		var plans string
		var start string
		var end string
		var created string

		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&start,
			&end,
			&plans,
			&created,
			&s.CreatedBy,
		)

		if err != nil {
			logger.WithField("location", "ListDiscount").Error(err)

			continue
		}

		if err := json.Unmarshal([]byte(plans), &s.Plans); err != nil {
			continue
		}

		s.Start = util.ISO8601UTC.FromDatetime(start, nil)
		s.End = util.ISO8601UTC.FromDatetime(end, nil)
		s.CreatedAt = util.ISO8601UTC.FromDatetime(created, nil)

		schs = append(schs, s)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "ListDiscounts").Error(err)

		return schs, err
	}

	return schs, nil
}

// DeleteSchedule delete a schedule row.
func (env Env) DeleteSchedule(id int64) error {
	query := `
	DELETE FROM premium.discount_schedule
	WHERE id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, id)

	if err != nil {
		logger.WithField("location", "DeleteDiscount").Error(err)
		return err
	}

	return nil
}
