package repository

import (
	log "github.com/sirupsen/logrus"
)

var logger = log.WithField("package", "model")

const (
	stmtStaffAccount = `
	SELECT id AS id,
		IFNULL(email, '') AS email,
		user_name AS userName,
		is_active AS isActive,
		display_name AS displayName,
		department AS department,
		group_memberships
	FROM backyard.staff`

	stmtStaffProfile = `
	SELECT id AS id,
		IFNULL(email, '') AS email,
		user_name AS userName,
		is_active AS isActive,
		display_name AS displayName,
		department AS department,
		group_memberships AS groups,
	    created_utc AS createdAt,
		deactivated_utc AS deactivatedAt,
		updated_utc AS updatedAt,
		last_login_utc AS lastLoginAt,
		INET6_NTOA(staff.last_login_ip) AS lastLoginIp
  	FROM backyard.staff`

	stmtUser = `
	SELECT user_id AS id,
		wx_union_id AS unionId,
		email AS email,
		user_name AS userName,
	    is_vip AS isVip
	FROM cmstmp01.userinfo`

	stmtOrder = `
	SELECT trade_no AS orderId,
		user_id AS userId,
		tier_to_buy AS tier,
		billing_cycle AS cycle,
	    trade_price AS listPrice,
	    trade_amount AS netPrice,
		payment_method AS payMethod,
		created_utc AS createdAt,
	    confirmed_utc AS confirmedAt,
		start_date AS startDate,
		end_date AS endDate,
		client_type AS clientType,
	    client_version AS clientVersion,
	    INET6_NTOA(user_ip_bin) AS userIp,
		user_agent AS userAgent
	FROM premium.ftc_trade`

	stmtFtcAccount = `
	SELECT u.user_id AS id,
		u.wx_union_id AS uUnionId,
		u.email AS email,
		u.user_name AS userName,
	    u.is_vip AS isVip,
	    u.mobile_phone_no AS mobile,
	    w.nickname AS nickName,
		IFNULL(v.vip_type, 0) AS vipType,
		IFNULL(v.expire_time, 0) AS expireTime,
		v.member_tier AS memberTier,
		v.billing_cycle AS billingCyce,
		v.expire_date AS expireDate,
		u.created_utc AS createdAt,
		u.updated_utc AS updatedAt
	FROM cmstmp01.userinfo AS u
		LEFT JOIN premium.ftc_vip AS v
		ON u.user_id = v.vip_id
		LEFT JOIN user_db.wechat_userinfo AS w
		ON u.wx_union_id = w.union_id
	WHERE u.%s = ?
	LIMIT 1`

	stmtWxAccount = `
	SELECT IFNULL(u.user_id, '') AS id,
		w.union_id AS wUnionId,
		IFNULL(u.email, '') AS email,
		u.user_name AS userName,
		IFNULL(u.is_vip, 0) AS isVip,
		u.mobile_phone_no AS mobile,
		w.nickname AS nickname,
		IFNULL(v.vip_type, 0) AS vipType,
		IFNULL(v.expire_time, 0) AS expireTime,
		v.member_tier AS memberTier,
		v.billing_cycle AS billingCyce,
		v.expire_date AS expireDate,
		w.created_utc AS createdAt,
		w.updated_utc AS updatedAt
	FROM user_db.wechat_userinfo AS w
		LEFT JOIN premium.ftc_vip AS v
		ON w.union_id = v.vip_id_alias
		LEFT JOIN cmstmp01.userinfo AS u
		ON w.union_id = u.wx_union_id
	WHERE w.union_id = ?
	LIMIT 1`

	stmtWxUser = `
	SELECT union_id AS unionId,
		nickname,
		avatar_url AS avatarUrl,
		gender,
		country,
		province,
		city,
		IFNULL(privilege, '') AS privilege,
	    created_utc AS createdAt,
	    updated_utc AS updatedAt
	FROM user_db.wechat_userinfo`
)

type sqlUserCol int

const (
	sqlUserColID sqlUserCol = iota
	sqlUserColEmail
	sqlUserColUnionID
)

func (c sqlUserCol) String() string {
	names := [...]string{
		"user_id",
		"email",
		"union_id",
	}

	if c < sqlUserColID || c > sqlUserColUnionID {
		return ""
	}

	return names[c]
}

type table int

const (
	tableUser table = iota
	tableStaff
)

func (t table) colName() string {
	return "user_name"
}

func (t table) colID() string {
	switch t {
	case tableUser:
		return "user_id"
	case tableStaff:
		return "id"
	default:
		return ""
	}
}

func (t table) colEmail() string {
	return "email"
}
