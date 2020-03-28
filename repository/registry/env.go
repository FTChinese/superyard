package registry

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Env wraps db.
type Env struct {
	DB *sqlx.DB
}

var logger = logrus.WithField("package", "repository/registry")

const stmtSelectApp = `
SELECT app_name,
	slug_name,
	LOWER(HEX(client_id)) AS client_id,
	LOWER(HEX(client_secret)) AS client_secret,
	repo_url,
	description,
	homepage_url AS home_url,
	callback_url,
	is_active,
	created_utc AS created_at,
	updated_utc AS updated_at,
	owned_by
FROM oauth.app_registry`

const stmtSelectToken = `
SELECT k.id AS id,
	LOWER(HEX(k.access_token)) AS token,
	k.is_active AS is_active,
	k.expires_in AS expires_in,
	k.usage_type AS usage_type,
	k.client_id AS client_id,
	k.description AS description,
	k.created_by AS created_by,
	k.created_utc AS created_at,
	k.updated_utc AS updated_at,
	k.last_used_utc AS last_used_at
FROM oauth.access AS k`
