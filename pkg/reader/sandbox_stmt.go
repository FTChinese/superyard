package reader

// StmtInsertSandbox records which account is sandbox and store the password as clear text.
const StmtInsertSandbox = `
INSERT INTO user_db.sandbox_account
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

const sandboxUserFrom = `
FROM user_db.sandbox_account AS s
	LEFT JOIN cmstmp01.userinfo AS u
	ON s.ftc_id = u.user_id
`

// StmtListSandboxUsers retrieves a list of FtcAccount.
const StmtListSandboxUsers = colsFtcAccount + sandboxUserFrom + `
WHERE u.user_id IS NOT NULL
ORDER BY u.created_UTC DESC`

const StmtSandboxJoinedAccount = colsJoinedAccount + `,
s.clear_password AS password,
s.created_by AS created_by
` + sandboxUserFrom + `
	LEFT JOIN user_db.wechat_userinfo AS w 
	ON u.wx_union_id = w.union_id
WHERE s.ftc_id = ?
	AND u.user_id IS NOT NULL
LIMIT 1`

const StmtSandboxExists = `
SELECT EXISTS(
	SELECT *
	FROM user_db.sandbox_account
	WHERE ftc_id = ?
) AS sandboxFound`

const StmtUpdateClearPassword = `
UPDATE user_db.sandbox_account
SET clear_password = :password
WHERE ftc_id = :ftc_id
LIMIT 1`

const StmtUpdatePassword = `
UPDATE cmstmp01.userinfo
SET password := MD5(:password),
	updated_utc := UTC_TIMESTAMP()
WHERE user_id = :ftc_id
LIMIT 1`
