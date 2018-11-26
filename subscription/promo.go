package subscription

import (
	"encoding/json"
	"fmt"

	"gitlab.com/ftchinese/backyard-api/util"
)

// Promotion contains all data for a promotion campaign.
// A pomotion compaign is divided into three steps:
// Schedule the time when it will start and end;
// Pricing plans for each products;
// Banner content used for pomotion.
// It is created not in one shot, but step by step.
// When retrieving, all data are retrieved together.
// When deleting, everything is deleted for a promotion record.
// When updating, schedule and plans parts are not allowed to edit;
// but banner content is editable.Promotion
type Promotion struct {
	Schedule
	Plans  map[string]Plan `json:"plans"`
	Banner Banner          `json:"banner"`
}

// RetrievePromo loads a promotion schedule record.
func (env Env) RetrievePromo(id int64) (Promotion, error) {
	query := fmt.Sprintf(`
	%s
	WHERE id = ?
	LIMIT 1`, stmtPromo)

	var p Promotion
	var startUtc string
	var endUtc string
	var plans string
	var banner string
	var createdUtc string

	err := env.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&startUtc,
		&endUtc,
		&plans,
		&banner,
		&p.CreatedAt,
		&p.CreatedBy,
	)

	if err != nil {
		logger.WithField("location", "RetrievePromo").Error(err)
		return p, err
	}

	if err := json.Unmarshal([]byte(plans), &p.Plans); err != nil {
		return p, err
	}

	if err := json.Unmarshal([]byte(banner), &p.Banner); err != nil {
		return p, err
	}

	p.Start = util.ISO8601UTC.FromDatetime(startUtc, nil)
	p.End = util.ISO8601UTC.FromDatetime(endUtc, nil)
	p.CreatedAt = util.ISO8601UTC.FromDatetime(createdUtc, nil)

	return p, nil
}

// ListPromo retrieves a list of promotion schedules by page.
func (env Env) ListPromo(page, rowCount int64) ([]Promotion, error) {
	offset := (page - 1) * rowCount

	query := fmt.Sprintf(`
	%s
	ORDER BY id DESC
	LIMIT ? OFFSET ?`, stmtPromo)

	rows, err := env.DB.Query(query, rowCount, offset)

	if err != nil {
		logger.WithField("location", "ListPromo").Error(err)

		return nil, err
	}

	defer rows.Close()

	promos := make([]Promotion, 0)

	for rows.Next() {
		var p Promotion
		var startUtc string
		var endUtc string
		var plans string
		var banner string
		var createdUtc string

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&startUtc,
			&endUtc,
			&plans,
			&banner,
			&p.CreatedAt,
			&p.CreatedBy,
		)

		if err != nil {
			logger.WithField("location", "ListPromo").Error(err)
			continue
		}

		if err := json.Unmarshal([]byte(plans), &p.Plans); err != nil {
			continue
		}

		if err := json.Unmarshal([]byte(banner), &p.Banner); err != nil {
			continue
		}

		p.Start = util.ISO8601UTC.FromDatetime(startUtc, nil)
		p.End = util.ISO8601UTC.FromDatetime(endUtc, nil)
		p.CreatedAt = util.ISO8601UTC.FromDatetime(createdUtc, nil)

		promos = append(promos, p)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "ListPromo").Error(err)

		return promos, err
	}

	return promos, nil
}

// EnablePromo turn a promotion record to enabled or disabled
func (env Env) EnablePromo(id int64, isEnabled bool) error {
	query := `
	UPDATE premium.promotion_schedule
	SET is_enabled = ?
	WHERE id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, isEnabled, id)

	if err != nil {
		logger.WithField("location", "DeletePromo").Error(err)
		return err
	}

	return nil
}
