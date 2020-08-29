package staff

const StmtLogin = StmtAccountCols + `
FROM backyard.staff AS s
WHERE (s.user_name, s.password) = (?, UNHEX(MD5(?)))
	AND s.is_active = 1`

const StmtUpdateLastLogin = `
UPDATE backyard.staff
SET last_login_utc = UTC_TIMESTAMP(),
	last_login_ip = IFNULL(INET6_ATON(?), last_login_ip)
WHERE user_name = ?
LIMIT 1`
