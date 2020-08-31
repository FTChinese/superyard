package reader

const ftcAccountCols = `
SELECT u.user_id AS ftc_id,
	w.union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name
`

// StmtFtcAccount retrieves ftc-only account by user_id.
const StmtFtcAccount = ftcAccountCols + `
FROM cmstmp01.userinfo AS u
WHERE u.user_id = ?
LIMIT 1`

const joinedAccountCols = ftcAccountCols + `,
	w.nickname AS wx_nickname,
	w.avatar_url AS wx_avatar_url,
	IFNULL(u.is_vip, FALSE) AS is_vip
`

const selectJoinedAccountByFtc = joinedAccountCols + `
FROM cmstmp01.userinfo AS u
	LEFT JOIN user_db.wechat_userinfo AS w 
	ON u.wx_union_id = w.union_id
`

// StmtJoinedAccountByFtcID select both ftc and wx account
// columns by ftc id. The wx columns might be zero values.
const StmtJoinedAccountByFtcID = selectJoinedAccountByFtc + `
WHERE u.user_id = ?
LIMIT 1`

const selectJoinedAccountByWx = joinedAccountCols + `
FROM user_db.wechat_userinfo AS w
	LEFT JOIN cmstmp01.userinfo AS u
	ON w.union_id = u.wx_union_id`

const StmtJoinedAccountByWxID = selectJoinedAccountByWx + `
WHERE w.union_id = ?
LIMIT 1`

// StmtSearchJoinedAccountByEmail retrieves FtcAccount by email.
const StmtSearchJoinedAccountByEmail = selectJoinedAccountByFtc + `
WHERE u.email LIKE ?
ORDER BY email ASC
LIMIT ? OFFSET ?`

const StmtSearchJoinedAccountByWxName = selectJoinedAccountByWx + `
WHERE w.nickname LIKE ?
ORDER BY nickname ASC
LIMIT ? OFFSET ?`
