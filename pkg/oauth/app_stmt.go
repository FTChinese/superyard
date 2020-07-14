package oauth

const StmtInsertApp = `
INSERT INTO oauth.app_registry
SET app_name = :app_name,
	slug_name = :slug_name,
	client_id = UNHEX(:client_id),
	client_secret = UNHEX(:client_secret),
	repo_url = :repo_url,
	description = :description,
	homepage_url = :home_url,
	callback_url = :callback_url,
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP(),
	owned_by = :owned_by`

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

const StmtApp = appCols + `
WHERE client_id = UNHEX(?)
LIMIT 1`

const StmtUpdateApp = `
UPDATE oauth.app_registry
SET app_name = :name,
	slug_name = :slug,
	repo_url = :repo_url,
	description = :description,
	homepage_url = :home_url,
	callback_url = :callback_url,
	updated_utc = UTC_TIMESTAMP()
WHERE client_id = UNHEX(:client_id)
	AND owned_by = :owned_by
	AND is_active = 1
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
