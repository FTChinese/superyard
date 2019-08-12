package paywall

import (
	"encoding/json"
	"fmt"
	"github.com/FTChinese/go-rest"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"gitlab.com/ftchinese/backyard-api/models/subs"
)

// PromoEnv is used to manage promotion schedule.
type PromoEnv struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository.paywall")

// NewSchedule saves a new promotion schedule.
// Return the inserted row's id so that client knows which row to update in the following step.
func (env PromoEnv) NewSchedule(s subs.Schedule, creator string) (int64, error) {
	query := `
	INSERT INTO premium.promotion_schedule
	SET name = ?,
		description = ?,
		start_utc = ?,
		end_utc = ?,
		created_by = ?,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	result, err := env.DB.Exec(query,
		s.Name,
		s.Description,
		s.StartAt,
		s.EndAt,
		creator,
	)

	if err != nil {
		logger.WithField("trace", "NewSchedule").Error(err)
		return -1, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		logger.WithField("trace", "NewSchedule").Error(err)
		return -1, err
	}

	return id, nil
}

// SavePlans set the pricing plans of a promotion schedule.
func (env PromoEnv) SavePlans(id int64, plans subs.Pricing) error {
	query := `
	UPDATE premium.promotion_schedule
	SET plans = ?,
		updated_utc = UTC_TIMESTAMP()
	WHERE id = ?
	LIMIT 1`

	p, err := json.Marshal(plans)

	if err != nil {
		logger.WithField("trace", "SetPlans").Error(err)
		return err
	}

	_, err = env.DB.Exec(query, string(p), id)

	if err != nil {
		logger.WithField("trace", "SavePlans").Error(err)
		return err
	}

	return nil
}

// SaveBanner sets the banner content for a promotion.
// It is also used to edit banner content.
func (env PromoEnv) SaveBanner(id int64, banner subs.Banner) error {
	query := `
	UPDATE premium.promotion_schedule
	SET banner = ?,
		updated_utc = UTC_TIMESTAMP()
	WHERE id = ?
	LIMIT 1`

	b, err := json.Marshal(banner)

	if err != nil {
		logger.WithField("trace", "SaveBanner").Error(err)

		return err
	}

	_, err = env.DB.Exec(query, string(b), id)

	if err != nil {
		logger.WithField("trace", "SaveBanner").Error(err)
		return err
	}

	return nil
}

// ListPromos retrieve a list of promotion schedules by page.
func (env PromoEnv) ListPromos(p gorest.Pagination) ([]subs.Promotion, error) {

	query := fmt.Sprintf(`
	%s
	WHERE is_enabled = 1
	ORDER BY id DESC
	LIMIT ? OFFSET ?`, stmtPromo)

	rows, err := env.DB.Query(query, p.Limit, p.Offset())

	if err != nil {
		logger.WithField("location", "ListPromo").Error(err)

		return nil, err
	}

	defer rows.Close()

	promos := make([]subs.Promotion, 0)

	for rows.Next() {
		var p subs.Promotion
		var plans string
		var banner string

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.StartAt,
			&p.EndAt,
			&plans,
			&banner,
			&p.IsEnabled,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CreatedBy,
		)

		if err != nil {
			logger.WithField("trace", "ListPromos").Error(err)
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

		promos = append(promos, p)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("trace", "ListPromos").Error(err)

		return promos, err
	}

	return promos, nil
}

// LoadPromo loads a promotion schedule record.
func (env PromoEnv) LoadPromo(id int64) (subs.Promotion, error) {
	query := fmt.Sprintf(`
	%s
	WHERE id = ?
		AND is_enabled = 1
	LIMIT 1`, stmtPromo)

	var p subs.Promotion
	var plans string
	var banner string

	err := env.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.StartAt,
		&p.EndAt,
		&plans,
		&banner,
		&p.IsEnabled,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.CreatedBy,
	)

	if err != nil {
		logger.WithField("trae", "LoadPromo").Error(err)
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

	return p, nil
}

// DisablePromo turn a promotion record to enabled or disabled.
func (env PromoEnv) DisablePromo(id int64) error {
	query := `
	UPDATE premium.promotion_schedule
	SET is_enabled = 0
	WHERE id = ?
	LIMIT 1`

	_, err := env.DB.Exec(query, id)

	if err != nil {
		logger.WithField("trace", "DeletePromo").Error(err)
		return err
	}

	return nil
}
