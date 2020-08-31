package reader

// StmtInsertSandbox records which account is sandbox and store the password as clear text.
const StmtInsertSandbox = `
INSERT INTO user_db.sandbox
SET ftc_id = :ftc_id,
	clear_password = :password,
	created_by = :created_by`

// StmtCreateReader insert the sandbox account.
const StmtCreateReader = `
INSERT INTO cmstmp01.userinfo
SET user_id = :ftc_id,
	email = :email,
	password = MD5(:password),
	created_utc = UTC_TIMESTAMP(),
	updated_utc = UTC_TIMESTAMP()`

const colSandboxUser = `
SELECT s.ftc_id AS ftc_id,
	u.wx_union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name
	s.created_by AS created_by,
	u.created_utc AS created_utc,
	u.updated_utc AS updated_utc
`

const sandboxUserFrom = `
FROM user_db.sandbox_account AS s
	LEFT JOIN cmstmp01.userinfo AS u
	ON s.ftc_id = u.user_id
`

const StmtSandboxUser = colSandboxUser + `,
s.clear_password AS password
` + sandboxUserFrom + `
WHERE s.ftc_id = ?
LIMIT 1`

const StmtListSandboxUsers = colSandboxUser + sandboxUserFrom + `
ORDER BY u.created_UTC DESC`

const StmtSandboxExists = `
SELECT EXISTS(
	SELECT *
	FROM user_db.sandbox_account
	WHERE ftc_id = ?
) AS sandboxFound`

const StmtUpdateClearPassword = `
UPDATE user_db.sandbox
SET password = :password
WHERE user_id := :ftc_id
LIMIT 1`

const StmtUpdatePassword = `
UPDATE cmstmp01.userinfo
SET password := MD5(:password),
	updated_utc := UTC_TIMESTAMP()
WHERE user_id := :ftc_id
LIMIT 1`
