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
	Plans     map[string]Plan `json:"plans"`
	Banner    *Banner         `json:"banner"`
	IsEnabled bool            `json:"isEnabled"`
	CreatedAt string          `json:"createdAt"`
	UpdatedAt string          `json:"updatedAt"`
	CreatedBy string          `json:"createdBy"`
}

// RetrievePromo loads a promotion schedule record.
func (env Env) RetrievePromo(id int64) (Promotion, error) {
	query := fmt.Sprintf(`
	%s
	WHERE id = ?
	LIMIT 1`, stmtPromo)

	var p Promotion
	var plans string
	var banner string

	err := env.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Start,
		&p.End,
		&plans,
		&banner,
		&p.IsEnabled,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.CreatedBy,
	)

	if err != nil {
		logger.WithField("location", "RetrievePromo").Error(err)
		return p, err
	}

	// Scanning a nullable JSON column is quite complicated.
	if plans != "" {
		if err := json.Unmarshal([]byte(plans), &p.Plans); err != nil {
			return p, err
		}
	}

	if banner != "" {
		if err := json.Unmarshal([]byte(banner), &p.Banner); err != nil {
			return p, err
		}
	}

	p.Start = util.ISO8601UTC.FromDatetime(p.Start, nil)
	p.End = util.ISO8601UTC.FromDatetime(p.End, nil)
	p.CreatedAt = util.ISO8601UTC.FromDatetime(p.CreatedAt, nil)
	p.UpdatedAt = util.ISO8601UTC.FromDatetime(p.UpdatedAt, nil)

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
		var plans string
		var banner string

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Start,
			&p.End,
			&plans,
			&banner,
			&p.IsEnabled,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CreatedBy,
		)

		if err != nil {
			logger.WithField("location", "ListPromo").Error(err)
			continue
		}

		if plans != "" {
			if err := json.Unmarshal([]byte(plans), &p.Plans); err != nil {
				continue
			}
		}

		if banner != "" {
			if err := json.Unmarshal([]byte(banner), &p.Banner); err != nil {
				continue
			}
		}

		p.Start = util.ISO8601UTC.FromDatetime(p.Start, nil)
		p.End = util.ISO8601UTC.FromDatetime(p.End, nil)
		p.CreatedAt = util.ISO8601UTC.FromDatetime(p.CreatedAt, nil)
		p.UpdatedAt = util.ISO8601UTC.FromDatetime(p.UpdatedAt, nil)

		promos = append(promos, p)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "ListPromo").Error(err)

		return promos, err
	}

	return promos, nil
}

// DisablePromo turn a promotion record to enabled or disabled.
func (env Env) DisablePromo(id int64) error {
	query := `
	UPDATE premium.promotion_schedule
	SET is_enabled = 0
	WHERE id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, id)

	if err != nil {
		logger.WithField("location", "DeletePromo").Error(err)
		return err
	}

	return nil
}
