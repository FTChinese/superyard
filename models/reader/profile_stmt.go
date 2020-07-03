package reader

const StmtFtcProfile = `
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

const StmtWxProfile = `
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
