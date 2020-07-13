package staff

const stmtProfileCols = StmtAccountCols + `,
	s.created_utc 				AS created_at,
	s.deactivated_utc 			AS deactivated_at,
	s.updated_utc 				AS updated_at,
	s.last_login_utc 			AS last_login_at,
	INET6_NTOA(s.last_login_ip) AS last_login_ip
FROM backyard.staff AS s`

const StmtActiveProfile = stmtProfileCols + `
WHERE s.staff_id = ?
	AND s.is_active = 1
LIMIT 1`

const StmtProfile = stmtProfileCols + `
WHERE s.staff_id = ?
LIMIT 1`
