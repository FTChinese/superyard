package subscription

import (
	"encoding/json"
	"fmt"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Promo contains all data for a promotion campaign.
type Promotion struct {
	Schedule
	Plans  map[string]Plan `json:"plans"`
	Banner Banner          `json:"banner"`
}

// RetrieveSchedule loads a schedule from database.
func (env Env) RetrieveSchedule(id int64) (Promotion, error) {
	query := fmt.Sprintf(`
	%s
	WHERE id = ?
	LIMIT 1`, stmtDiscount)

	var s Promotion
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
func (env Env) ListSchedules(page, rowCount int64) ([]Promotion, error) {
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

	schs := make([]Promotion, 0)

	for rows.Next() {
		var s Promotion
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
