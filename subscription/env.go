package subscription

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const (
	stmtDiscount = `SELECT 
		id AS id,
		name AS name,
		description AS description,
		start_utc AS start,
		end_utc AS end,
		plans AS plans,
		created_utc AS createdUtc,
		created_by AS createdBy
	FROM premium.discount_schedule`

	stmtPromo = `SELECT
		id AS id,
		name AS name,
		description AS description,
		start_utc AS startUtc,
		end_utc AS endUtc,
		plans AS plans,
		banner AS banner,
		created_utc AS createdUtc,
		created_by AS createdBy
	FROM premium.promotion_schedule`
)

var logger = log.WithField("package", "subscription")

// Env wraps database connection.
type Env struct {
	DB *sql.DB
}
