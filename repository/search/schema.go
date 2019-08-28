package search

const (
	stmtSearchFtc = `
	SELECT user_id AS ftc_id,
		email
	FROM cmstmp01.userinfo
	WHERE email = ?
	LIMIT 1`

	stmtSearchWx = `
	SELECT w.union_id,
		w.nickname
	FROM user_db.wechat_userinfo
	WHERE nickname LIKE ?
	ORDER BY nickname ASC
	LIMIT ? OFFSET ?`
)
