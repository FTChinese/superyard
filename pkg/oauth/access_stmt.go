package oauth

const StmtInsertToken = `
INSERT INTO oauth.access
SET access_token = UNHEX(:token),
	is_active = :is_active,
	expires_in = :expires_in,
	usage_type = :usage_type,
	description = :description,
	created_by = :created_by,
	client_id = UNHEX(:client_id),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const accessTokenCols = `
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

const StmtListAppKeys = accessTokenCols + `
WHERE k.is_active = 1
	AND k.client_id = UNHEX(?)
	AND k.usage_type = 'app'
ORDER BY k.created_utc DESC
LIMIT ? OFFSET ?`

const StmtListPersonalKeys = accessTokenCols + `
WHERE k.is_active = 1
	AND k.created_by = ?
	AND k.usage_type = 'personal'
ORDER BY k.created_utc DESC
LIMIT ? OFFSET ?`

const StmtRemoveToken = `
UPDATE oauth.access
	SET is_active = 0
WHERE id = :id
	AND created_by = :created_by
LIMIT 1`
