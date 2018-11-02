package subscription

import (
	"database/sql"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ftchinese/backyard-api/util"
)

const (
	stmtDiscount = `SELECT name AS name,
		description AS description,
		start_utc AS start,
		end_utc AS end,
		plans AS plans,
		created_utc
	FROM premium.discount_schedule`
)

var logger = log.WithField("package", "subscription")

// Env wraps database connection.
type Env struct {
	DB *sql.DB
}

// Plan contains details of subscription plan.
type Plan struct {
	Tier  string  `json:"tier"`
	Cycle string  `json:"cycle"`
	Price float64 `json:"price"`
	ID    int
	// For wxpay, this is used as `body` parameter;
	// For alipay, this is used as `subject` parameter.
	Description string `json:"description"`
	// For wxpay, this is used as `detail` parameter;
	// For alipay, this is used as `body` parameter.
	Message string `json:"message"`
}

// Discount represents a discount activity
type Discount struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Start       string          `json:"startAt"`
	End         string          `json:"endAt"`
	Plans       map[string]Plan `json:"plans"`
	CreatedAt   string          `json:"createdAt"`
}

// NewDiscount Create a new discount record.
func (env Env) NewDiscount(d Discount) error {
	query := `
	INSERT INTO premium.discount_schedule
	SET name = ?,
		description = ?,
		start_utc = ?,
		end_utc = ?,
		plans = ?`

	startUTC := util.SQLDatetimeUTC.FromISO8601(d.Start)
	endUTC := util.SQLDatetimeUTC.FromISO8601(d.End)
	plans, err := json.Marshal(d.Plans)

	if err != nil {
		return err
	}

	log.Println(plans)

	_, err = env.DB.Exec(query,
		d.Name,
		d.Description,
		startUTC,
		endUTC,
		string(plans),
	)

	if err != nil {
		return err
	}

	return nil
}

func (env Env) Retrieve(id int64) (Discount, error) {
	query := fmt.Sprintf(`
	%s
	WHERE id = ?
	LIMIT 1`, stmtDiscount)

	var d Discount
	var plans string
	var start string
	var end string
	var created string
	err := env.DB.QueryRow(query, id).Scan(
		&d.Name,
		&d.Description,
		&start,
		&end,
		&plans,
		&created,
	)

	if err != nil {
		return d, err
	}

	if err := json.Unmarshal([]byte(plans), &d.Plans); err != nil {
		return d, err
	}

	d.Start = util.ISO8601UTC.FromDatetime(start, nil)
	d.End = util.ISO8601UTC.FromDatetime(end, nil)
	d.CreatedAt = util.ISO8601UTC.FromDatetime(created, nil)

	return d, nil
}

func (env Env) ListDiscount(page uint, rowCount uint) ([]Discount, error) {
	offset := (page - 1) * rowCount

	query := fmt.Sprintf(`
	%s
	ORDER BY id DESC
	LIMIT ? OFFSET ?`, stmtDiscount)

	rows, err := env.DB.Query(query, rowCount, offset)

	if err != nil {
		logger.
			WithField("location", "ListDiscount").
			Error(err)

		return nil, err
	}
	defer rows.Close()

	discounts := make([]Discount, 0)

	for rows.Next() {
		var d Discount
		var plans string
		var start string
		var end string
		var created string

		err := rows.Scan(
			&d.Name,
			&d.Description,
			&start,
			&end,
			&plans,
			&created,
		)

		if err != nil {
			logger.WithField("location", "ListDiscount").Error(err)

			continue
		}

		if err := json.Unmarshal([]byte(plans), &d.Plans); err != nil {
			continue
		}

		d.Start = util.ISO8601UTC.FromDatetime(start, nil)
		d.End = util.ISO8601UTC.FromDatetime(end, nil)
		d.CreatedAt = util.ISO8601UTC.FromDatetime(created, nil)

		discounts = append(discounts, d)
	}

	if err := rows.Err(); err != nil {
		logger.WithField("location", "ListDiscounts").Error(err)

		return discounts, err
	}

	return discounts, nil
}

func (env Env) DeleteDiscount(id int64) error {
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
