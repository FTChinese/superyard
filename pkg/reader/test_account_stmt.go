package reader

// StmtInsertTestAccount records which account is sandbox and store the password as clear text.
const StmtInsertTestAccount = `
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

const StmtCreateProfile = `
INSERT INTO user_db.profile
SET user_id = :ftc_id`

const StmtDeleteTestUser = `
DELETE FROM user_db.sandbox_account
WHERE ftc_id = ?
LIMIT 1`

const StmtDeleteAccount = `
DELETE FROM cmstmp01.userinfo
WHERE user_id = ?
LIMIT 1`

const StmtDeleteProfile = `
DELETE FROm cmstmp01.userinfo
WHERE user_id = ?
LIMIT 1`

// StmtDeleteMember deletes the membership under a sandbox account.
// Never delete a real user's membership.
const StmtDeleteMember = `
DELETE FROM premium.ftc_vip
WHERE vip_id = ?
LIMIT 1`

const testUserFrom = `
FROM user_db.sandbox_account AS s
	LEFT JOIN cmstmp01.userinfo AS u
	ON s.ftc_id = u.user_id
`

const StmtCountTestUser = `
SELECT COUNT(*) AS row_count
FROM user_db.sandbox_account`

// StmtListTestUsers retrieves a list of FtcAccount.
const StmtListTestUsers = colsFtcAccount +
	testUserFrom + `
WHERE u.user_id IS NOT NULL
ORDER BY u.created_UTC DESC`

// StmtTestJoinedAccount is similar to StmtJoinedAccountByFtcId with two extra columns.
const StmtTestJoinedAccount = colsJoinedAccount + `,
s.clear_password AS password,
s.created_by AS created_by
` + testUserFrom + `
	LEFT JOIN user_db.wechat_userinfo AS w 
	ON u.wx_union_id = w.union_id
WHERE s.ftc_id = ?
	AND u.user_id IS NOT NULL
LIMIT 1`

const StmtTestUserExists = `
SELECT EXISTS(
	SELECT *
	FROM user_db.sandbox_account
	WHERE ftc_id = ?
) AS sandboxFound`

const StmtUpdateTestUserPassword = `
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
