package search

const (
	stmtSearchFtc = `
	SELECT user_id AS ftc_id,
		email,
		is_vip
	FROM cmstmp01.userinfo
	WHERE email = ?
	LIMIT 1`

	stmtSearchWx = `
	SELECT union_id,
		nickname
	FROM user_db.wechat_userinfo
	WHERE nickname LIKE ?
	ORDER BY nickname ASC
	LIMIT ? OFFSET ?`
)
