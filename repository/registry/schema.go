package registry

const (
	stmtCreateApp = `
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

	stmtSelectApp = `
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

	// Use the following to enforce access control:
	//	WHERE client_id = UNHEX(?)
	//      	AND owned_by = ?
	//      	AND is_active = 1
	stmtRemoveApp = `
	UPDATE oauth.app_registry
      	SET is_active = 0
	WHERE client_id = UNHEX(:client_id)
      	AND owned_by = :owned_by
      	AND is_active = 1
	LIMIT 1`

	// Remove all keys belong to an app when deleting an app.
	stmtRemoveAppKeys = `
	UPDATE oauth.access
     	SET is_active = 0
    WHERE client_id = UNHEX(:client_id)
		AND usage_type = 'app'`

	//stmtTransferApp = `
	//UPDATE oauth.app_registry
	//	SET owned_by = ?,
	//		updated_utc = UTC_TIMESTAMP()
	//WHERE slug_name = ?
	//	AND owned_by = ?
	//  	AND is_active = 1
	//LIMIT 1`

	stmtInsertToken = `
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

	stmtSelectToken = `
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

	// Retrieve keys owned by an app.
	stmtAppKeys = stmtSelectToken + `
	WHERE k.is_active = 1
		AND k.client_id = UNHEX(?)
		AND k.usage_type = 'app'
	ORDER BY k.created_utc DESC
	LIMIT ? OFFSET ?`

	// Retrieve a staff's personal keys.
	stmtPersonalKeys = stmtSelectToken + `
	WHERE k.is_active = 1
		AND k.created_by = ?
		AND k.usage_type = 'personal'
	ORDER BY k.created_utc DESC
	LIMIT ? OFFSET ?`

	// Deactivate all personal keys created by a user.
	stmtRemovePersonalKeys = `
	UPDATE oauth.access
		SET is_active = 0
	WHERE created_by = ?
		AND usage_type = 'personal'`

	// Deactivate a key by whoever created it.
	stmtRemoveKey = `
	UPDATE oauth.access
		SET is_active = 0
	WHERE id = :id
		AND created_by = :created_by
	LIMIT 1`
)
