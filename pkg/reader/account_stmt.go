package reader

const colsFtcAccount = `
SELECT u.user_id AS ftc_id,
	u.wx_union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name,
	u.created_utc AS created_utc,
	u.updated_utc AS updated_utc,
	IFNULL(u.is_vip, FALSE) AS is_vip
`

const stmtSelectFtcAccount = colsFtcAccount + `
FROM cmstmp01.userinfo AS u
`

// StmtFtcAccount retrieves ftc-only account by user_id.
const StmtFtcAccount = stmtSelectFtcAccount + `
WHERE u.user_id = ?
LIMIT 1`

const StmtFindFtcAccount = stmtSelectFtcAccount + `
WHERE ? IN (u.email, u.user_name)
LIMIT 1`

const StmtCountVIP = `
SELECT COUNT(*) AS row_count
FROM cmstmp01.userinfo
WHERE is_vip = TRUE`

const StmtListVIP = stmtSelectFtcAccount + `
WHERE u.is_vip = TRUE
ORDER BY u.email ASC
LIMIT ? OFFSET ?`

const StmtSetVIP = `
UPDATE cmstmp01.userinfo
SET is_vip = :is_vip
WHERE user_id = :ftc_id
LIMIT 1`

const colsJoinedAccount = colsFtcAccount + `,
	w.nickname AS wx_nickname,
	w.avatar_url AS wx_avatar_url
`

const selectJoinedAccountByFtc = colsJoinedAccount + `
FROM cmstmp01.userinfo AS u
	LEFT JOIN user_db.wechat_userinfo AS w 
	ON u.wx_union_id = w.union_id
`

const selectJoinedAccountByWx = colsJoinedAccount + `
FROM user_db.wechat_userinfo AS w
	LEFT JOIN cmstmp01.userinfo AS u
	ON w.union_id = u.wx_union_id`

// StmtSearchJoinedAccountByEmail retrieves FtcAccount by email.
const StmtSearchJoinedAccountByEmail = selectJoinedAccountByFtc + `
WHERE u.email LIKE ?
ORDER BY email ASC
LIMIT ? OFFSET ?`

const StmtSearchJoinedAccountByWxName = selectJoinedAccountByWx + `
WHERE w.nickname LIKE ?
ORDER BY nickname ASC
LIMIT ? OFFSET ?`
