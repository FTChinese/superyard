package registry

const (
	stmtFTCApp = `
	SELECT id AS id,
		app_name AS appName,
    	slug_name AS slugName,
    	LOWER(HEX(client_id)) AS clientId,
    	LOWER(HEX(client_secret)) AS clientSecret,
    	repo_url AS repoUrl,
    	description AS description,
    	homepage_url AS homeUrl,
		callback_url AS callbackUrl,
		is_active AS isActive,
		created_utc AS createdAt,
		updated_utc AS updatedAt,
    	owned_by AS ownedBy
	FROM oauth.app_registry`

	stmtPersonalToken = `
	SELECT a.id AS id,
		LOWER(HEX(a.access_token)) AS token,
	    a.description AS description,
	    u.email AS ftcEmail,
	    a.created_by AS createdBy,
		a.created_utc AS createdAt,
		a.updated_utc AS updatedAt,
		a.last_used_utc AS lastUsedAt
	FROM oauth.access AS a
		LEFT JOIN cmstmp01.userinfo AS u
		ON a.myft_id = u.user_id
	WHERE a.is_active = 1
		AND a.created_by = ?
		AND a.client_id IS NULL`
)
