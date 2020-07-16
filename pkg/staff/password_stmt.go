package staff

const StmtInsertPwResetSession = `
INSERT INTO backyard.password_reset
SET token = UNHEX(:token),
	email = :email,
	created_utc = UTC_TIMESTAMP()`

// StmtPwResetSession retrieves a password reset session.
const StmtPwResetSession = `
SELECT LOWER(HEX(token)) AS token,
	email,
	is_used,
	expires_in,
	created_utc
FROM backyard.password_reset
WHERE token = UNHEX(?)
LIMIT 1`

// StmtAccountByResetToken retrieves an account for a password
// reset token. This used when user submitted new password together with the session token.
// We don't care whether the token is expires or not, as long as
// it is not used yet.
const StmtAccountByResetToken = StmtAccountCols + `
FROM backyard.password_reset AS r
	JOIN backyard.staff AS s
	ON r.email = s.email
WHERE r.token = UNHEX(?)
  AND r.is_used = 0
  AND s.is_active = 1
LIMIT 1`

const StmtDisableResetToken = `
UPDATE backyard.password_reset
SET is_used = 1
WHERE token = UNHEX(?)
LIMIT 1`

const StmtVerifyPassword = StmtAccountCols + `
FROM backyard.staff AS s
WHERE (s.staff_id, s.password) = (?, UNHEX(MD5(?)))
	AND s.is_active = 1
LIMIT 1`

const StmtUpdatePassword = `
UPDATE backyard.staff
SET password = UNHEX(MD5(:password)),
	updated_utc = UTC_TIMESTAMP()
WHERE user_name = :user_name
	AND is_active = 1
LIMIT 1`

const StmtUpdateLegacyPassword = `
UPDATE cmstmp01.managers
	SET password = MD5(:password)
WHERE username = :user_name
LIMIT 1`
