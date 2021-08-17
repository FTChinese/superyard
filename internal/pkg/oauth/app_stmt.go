package oauth

const upsertAppCols = `
app_name = :app_name,
slug_name = :slug_name,
repo_url = :repo_url,
description = :description,
homepage_url = :home_url,
callback_url = :callback_url,
updated_utc = UTC_TIMESTAMP()
`

const StmtInsertApp = `
INSERT INTO oauth.app_registry
SET client_id = UNHEX(:client_id),
	client_secret = UNHEX(:client_secret),
	created_utc = UTC_TIMESTAMP(),
	owned_by = :owned_by,` + upsertAppCols

const StmtUpdateApp = `
UPDATE oauth.app_registry
SET` + upsertAppCols + `
WHERE client_id = UNHEX(:client_id)
	AND is_active = 1
LIMIT 1`

const appCols = `
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

const StmtListApps = appCols + `
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

const StmtCountApp = `
SELECT COUNT(*) AS row_count
FROM oauth.app_registry`

const StmtApp = appCols + `
WHERE client_id = UNHEX(?)
LIMIT 1`

const StmtRemoveApp = `
UPDATE oauth.app_registry
	SET is_active = 0
WHERE client_id = UNHEX(?)
	AND is_active = 1
LIMIT 1`

const StmtRemoveAppKeys = `
UPDATE oauth.access
	SET is_active = 0
WHERE client_id = UNHEX(:client_id)
	AND usage_type = 'app'`
