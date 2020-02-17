package readers

const selectAccountByFtcID = `
SELECT u.user_id AS ftc_id,
	u.wx_union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name,
	w.nickname AS nickname,
	u.is_vip AS is_vip
FROM cmstmp01.userinfo AS u
	LEFT JOIN user_db.wechat_userinfo AS w 
	ON u.wx_union_id = w.union_id
WHERE u.user_id = ?
LIMIT 1`

const selectFtcIDByEmail = `
SELECT user_id AS ftc_id
FROM cmstmp01.userinfo
WHERE email = ?
LIMIT 1`

const selectWxAccountBase = `
SELECT u.user_id AS ftc_id,
	w.union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name,
	w.nickname AS nickname,
	u.is_vip AS is_vip
FROM user_db.wechat_userinfo AS w
	LEFT JOIN cmstmp01.userinfo AS u
	ON w.union_id = u.wx_union_id`

const selectAccountByWxID = selectWxAccountBase + `
WHERE w.union_id = ?
LIMIT 1`

const selectWxIDs = `
SELECT union_id
FROM user_db.wechat_userinfo
WHERE nickname LIKE ?
ORDER BY nickname ASC
LIMIT ? OFFSET ?`

const selectWxAccounts = selectWxAccountBase + `
WHERE FIND_IN_SET(w.union_id, ?) > 0
ORDER BY w.nickname ASC`

const selectFtcProfile = `
SELECT u.user_id AS ftc_id,
	u.wx_union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email,
	u.user_name,
	u.mobile_phone_no AS mobile,
	IFNULL(u.is_vip, 0) AS is_vip,
	u.gender,
	u.last_name,
	u.first_name,
	u.birthday,
	p.country AS country,
	p.province AS province,
	p.city AS city,
	p.district AS district,
	p.street AS street,
	p.postcode AS postcode,
	u.created_utc AS created_at,
	u.updated_utc AS updated_at
FROM cmstmp01.userinfo AS u
	LEFT JOIN user_db.profile AS p
	ON u.user_id = p.user_id
WHERE u.user_id = ?
LIMIT 1`

const selectWxProfile = `
SELECT union_id,
	nickname,
	avatar_url,
	gender,
	country,
	province,
	city,
	created_utc AS created_at,
	updated_utc AS updated_at
FROM user_db.wechat_userinfo
WHERE union_id = ?
LIMIT 1`

const (
	stmtSelectVIP = `
	SELECT user_id
		email,
		is_vip
	FROM cmstmp01.userinfo
	WHERE is_vip = 1
	LIMIT ? OFFSET ?`

	stmtGrantVIP = `
	UPDATE cmstmp01.userinfo
      SET is_vip = 1
    WHERE user_id = ?
	LIMIT 1`

	stmtRevokeVIP = `
	UPDATE cmstmp01.userinfo
      SET is_vip = 0
    WHERE user_id = ?
	LIMIT 1`
)
