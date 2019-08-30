package search

const (
	stmtSearchFtc = `
	SELECT user_id AS ftc_id,
		email
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

	sqlSearchStaff = `
	SELECT staff_id,
		IFNULL(email, '') AS email,
		user_name,
		is_active,
		display_name,
		department,
		group_memberships
	FROM backyard.staff`
)
