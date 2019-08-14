package registry

const (
	stmtCreateApp = `
	INSERT INTO oauth.app_registry
	SET app_name = ?,
		slug_name = ?,
        client_id = UNHEX(?),
        client_secret = UNHEX(?),
        repo_url = ?,
        description = ?,
        homepage_url = ?,
        callback_url = ?,
        created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP(),
		owned_by = ?`

	stmtSelectApp = `
	SELECT id,
		app_name,
    	slug_name,
    	LOWER(HEX(client_id)) AS client_id,
    	LOWER(HEX(client_secret)) AS client_secret,
    	repo_url,
    	description AS description,
    	homepage_url AS home_url,
		callback_url,
		is_active,
		created_utc AS created_at,
		updated_utc AS updated_at,
    	owned_by AS owned_by
	FROM oauth.app_registry`

	stmtListApps = stmtSelectApp + `
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	stmtApp = stmtSelectApp + `
	WHERE client_id = UNHEX(?)
	LIMIT 1`

	stmtUpdateApp = `
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

	stmtSearchApp = `
	SELECT LOWER(HEX(client_id)) AS client_id
	FROM oauth.app_registry
	WHERE slug_name = ?`

	// Use the following to enforce access control:
	//	WHERE client_id = UNHEX(?)
	//      	AND owned_by = ?
	//      	AND is_active = 1
	stmtRemoveApp = `
	UPDATE oauth.app_registry
      	SET is_active = 0
	WHERE client_id = UNHEX(?)
      	AND owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	stmtRemoveAppKeys = `
	UPDATE oauth.access
    	SET is_active = 0
    WHERE client_id = ?`

	stmtTransferApp = `
	UPDATE oauth.app_registry
    	SET owned_by = ?,
			updated_utc = UTC_TIMESTAMP()
	WHERE slug_name = ?
		AND owned_by = ?
      	AND is_active = 1
	LIMIT 1`

	stmtInsertToken = `
	INSERT INTO oauth.access
    SET access_token = UNHEX(:token),
		is_active = :is_active,
		expires_in = :expires_in,
		usage_type = :usage_type,
		ftc_id = :ftc_id,
		client_id = UNHEX(:client_id),
    	description = :description,
		created_by = :created_by,
		created_utc = UTC_TIMESTAMP(),
		updated_utc = UTC_TIMESTAMP()`

	stmtSelectToken = `
	SELECT k.id AS id,
		LOWER(HEX(k.access_token)) AS token,
		k.is_active,
		k.expires_in,
		k.usage_type AS usage,
		k.ftc_id,
		k.client_id,
		k.description,
		k.created_by,
		k.created_utc AS created_at,
		k.updated_utc AS updated_at,
		k.last_used_utc AS last_used_at
	FROM oauth.access AS k`

	// Retrieve keys owned by an app.
	stmtAppKeys = stmtSelectToken + `
	WHERE k.is_active = 1
		AND k.client_id = UNHEX(?)
		AND k.ftc_id IS NULL
	ORDER BY k.created_utc DESC
	LIMIT ? OFFSET ?`

	// Retrieve a staff's personal keys.
	stmtPersonalKeys = stmtSelectToken + `
	WHERE k.is_active = 1
		AND k.created_by = ?
		AND k.client_id IS NULL
	ORDER BY k.created_utc DESC
	LIMIT ? OFFSET ?`

	stmtRemoveAppKey = `
	UPDATE oauth.access
      SET is_active = 0
    WHERE access_token = ?
	  AND client_id = UNHEX(?)
	LIMIT 1`

	stmtRemovePersonalKey = `
	UPDATE oauth.access
		SET is_active = 0
	WHERE access_token = ?
		AND created_by = ?
	LIMIT 1`

	//stmtAppToken = `SELECT t.id AS id,
	//	LOWER(HEX(t.access_token)) AS token,
	//	t.description AS description,
	//	t.created_utc AS createdAt,
	//	t.updated_utc AS updatedAt,
	//	t.last_used_utc AS lastUsedAt
	//FROM oauth.access AS t
	//	JOIN oauth.app_registry AS a
	//	ON t.client_id = a.client_id
	//WHERE t.is_active = 1
	//	AND a.slug_name = ?
	//ORDER BY t.created_utc DESC
	//LIMIT ? OFFSET ?`

	//stmtPersonalToken = `
	//SELECT a.id AS id,
	//	LOWER(HEX(a.access_token)) AS token,
	//    a.description AS description,
	//    u.email AS ftcEmail,
	//    a.created_by AS createdBy,
	//	a.created_utc AS createdAt,
	//	a.updated_utc AS updatedAt,
	//	a.last_used_utc AS lastUsedAt
	//FROM oauth.access AS a
	//	LEFT JOIN cmstmp01.userinfo AS u
	//	ON a.myft_id = u.user_id
	//WHERE a.is_active = 1
	//	AND a.created_by = ?
	//	AND a.client_id IS NULL`
)
