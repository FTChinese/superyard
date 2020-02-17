package readers

const stmtActivity = `
SELECT ftc_id,
	platform,
	clint_version AS version,
	user_ip,
	user_agent,
	created_utc,
	source AS kind,
FROM user_db.client_footprint
WHERE ftc_id = ?
ORDER by created_utc DESC
LIMIT ? OFFSET ?`

const stmtWxLoginHistory = `
	SELECT union_id,
		open_id,
		app_id,
		client_type,
		client_version,
		INET6_NTOA(user_ip) AS user_ip,
		user_agent AS user_agent,
		created_utc AS created_at,
		updated_utc AS updated_at
	FROM user_db.wechat_access
	WHERE union_id = ?
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`
