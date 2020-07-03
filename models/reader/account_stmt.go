package reader

const ftcAccountCols = `
SELECT u.user_id AS ftc_id,
	w.wx_union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name
`

const ftcWxAccountCols = ftcAccountCols + `,
	w.nickname AS wx_nickname,
	w.avatar_url AS wx_avatar_url
	IFNULL(u.is_vip, FALSE) AS is_vip
`

const selectFtcAccount = ftcWxAccountCols + `
FROM cmstmp01.userinfo AS u
	LEFT JOIN user_db.wechat_userinfo AS w 
	ON u.wx_union_id = w.union_id
`
const StmtAccountByFtcID = selectFtcAccount + `
WHERE u.user_id = ?
LIMIT 1`

// StmtSearchFtcByEmail retrieves FtcAccount by email.
const StmtSearchFtcByEmail = selectFtcAccount + `
WHERE u.email LIKE ?
ORDER BY email ASC
LIMIT ? OFFSET ?`

const selectWxAccount = ftcWxAccountCols + `
FROM user_db.wechat_userinfo AS w
	LEFT JOIN cmstmp01.userinfo AS u
	ON w.union_id = u.wx_union_id`

const StmtAccountByWxID = selectWxAccount + `
WHERE w.union_id = ?
LIMIT 1`

const StmtSearchWxByName = selectWxAccount + `
WHERE w.nickname LIKE ?
ORDER BY nickname ASC
LIMIT ? OFFSET ?`
