package model

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

	stmtPromo = `
	SELECT id AS id,
		name AS name,
		description AS description,
		start_utc AS startUtc,
		end_utc AS endUtc,
		IFNULL(plans, '') AS plans,
		IFNULL(banner, '') AS banner,
		is_enabled AS isEnabled,
		created_utc AS createdUtc,
		updated_utc AS updatedUtc,
		created_by AS createdBy
	FROM premium.promotion_schedule`

	stmtFTCApp = `
	SELECT id AS id,
		app_name AS appName,
    	slug_name AS slugName,
    	LOWER(HEX(client_id)) AS clientId,
    	LOWER(HEX(client_secret)) AS clientSecret,
    	repo_url AS repoUrl,
    	description AS description,
    	homepage_url AS homeUrl,
		is_active AS isActive,
		created_utc AS createdAt,
		updated_utc AS updatedAt,
    	owned_by AS ownedBy
	FROM oauth.app_registry`

	stmtPersonalToken = `
	SELECT a.id AS id,
		LOWER(HEX(a.access_token)) AS token,
	    a.description AS description,
	    u.email AS ftcEmail,
	    a.created_by AS createdBy,
		a.created_utc AS createdAt,
		a.updated_utc AS updatedAt,
		a.last_used_utc AS lastUsedAt
	FROM oauth.access AS a
		LEFT JOIN cmstmp01.userinfo AS u
		ON a.myft_id = u.user_id
	WHERE a.is_active = 1
		AND a.created_by = ?
		AND a.client_id IS NULL`
)

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
