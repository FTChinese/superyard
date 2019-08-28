package customer

const stmtFtcJoinWx = `
SELECT u.user_id AS ftc_id,
	u.wx_union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name,
	w.nickname AS nickname
FROM cmstmp01.userinfo AS u
	LEFT JOIN user_db.wechat_userinfo AS w
	ON u.wx_union_id  = w.union_id
WHERE u.user_id = ?
LIMIT 1`

const stmtWxJoinFtc = `
SELECT u.user_id AS ftc_id,
	w.union_id AS union_id,
	u.stripe_customer_id AS stripe_id,
	u.email AS email,
	u.user_name AS user_name,
	w.nickname AS nickname
FROM user_db.wechat_userinfo AS w
	LEFT JOIN cmstmp01.userinfo AS u
	ON w.union_id = u.wx_union_id
WHERE w.union_id = ?
LIMIT 1`

const selectFtcProfile = `
SELECT user_id AS ftc_id,
	wx_union_id AS union_id,
	stripe_customer_id AS stripe_id,
	email,
	user_name,
	mobile_phone_no AS mobile,
	IFNULL(is_vip, 0) AS is_vip,
	gender,
	last_name,
	first_name,
	birthday,
	p.country AS country,
	p.province AS province,
	p.city AS city,
	p.district AS district,
	p.street AS street,
	p.postcode AS postcode,
	created_utc AS created_at,
	updated_utc AS updated_at
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
