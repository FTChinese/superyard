package subscription

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

const (
	stmtPromo = `SELECT
		id AS id,
		name AS name,
		description AS description,
		start_utc AS startUtc,
		end_utc AS endUtc,
		IFNULL(plans, '') AS plans,
		IFNULL(banner, '') AS banner,
		is_enabled AS isEnabled,
		created_utc AS createdUtc,
		updated_utc AS updatedUtc,
		created_by AS createdBy
	FROM premium.promotion_schedule`
)

var logger = log.WithField("package", "subscription")

// Env wraps database connection.
type Env struct {
	DB *sql.DB
}
