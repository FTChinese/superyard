package model

const (
	stmtStaffAccount = `
	SELECT id AS id,
		username AS userName,
		IFNULL(email, '') AS email,
		display_name AS displayName,
		department AS department,
		group_memberships AS groups
	FROM backyard.staff`

	stmtStaffProfile = `
	SELECT id AS id,
		IFNULL(email, '') AS email,
	    username AS userName,
		display_name AS displayName,
		department AS department,
		group_memberships AS groups,
	    is_active AS isActive,
	    created_utc AS createdAt,
		deactivated_utc AS deactivatedAt,
		updated_utc AS updatedAt,
		last_login_utc AS lastLoginAt,
		INET6_NTOA(staff.last_login_ip) AS lastLoginIp
  	FROM backyard.staff`

	stmtMyft = `
	SELECT user_id AS id,
		email AS email,
	    is_vip AS isVip
	FROM cmstmp01.userinfo`

	stmtOrder = `
	SELECT trade_no AS orderId,
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
	    INET6_NTOA(user_ip_bin) AS userIp
	FROM premium.ftc_trade`

	stmtPromo = `SELECT
		id AS id,
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
)

type table int

const (
	tableUser table = iota
	tableStaff
)

func (t table) colName() string  {
	switch t {
	case tableUser:
		return "username"
	case tableStaff:
		return "user_name"
	default:
		return ""
	}
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