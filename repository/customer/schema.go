package customer

type MemberColumn string

const (
	MemberColumnFtcID MemberColumn = "vip_id"
	MemberColumnWxID               = "vip_id_alias"
)

const (

	// Select an FTC user by either user_id
	// or email column
	stmtFtcInfo = `
	SELECT user_id AS ftc_id,
		wx_union_id AS union_id,
		stripe_customer_id AS stripe_id,
		email,
		user_name,
	    IFNULL(is_vip, 0) AS is_vip,
	    mobile_phone_no AS mobile,
		created_utc AS created_at,
		updated_utc AS updated_at
	FROM cmstmp01.userinfo
	WHERE user_id = ?
	LIMIT 1`

	// Select a wechat account either by
	// union_id or nickname
	stmtWxAccount = `
	SELECT w.union_id AS union_id,
		w.nickname AS nickname,
	    w.created_utc AS created_at,
	    w.updated_utc AS updated_at,
		u.user_id AS ftc_id
	FROM user_db.wechat_userinfo AS w
		LEFT JOIN cmstmp01.userinfo AS u
		ON w.union_id = u.wx_union_id
	WHERE w.union_id = ?
	LIMIT 1`

	// Select membership by either vip_id
	// (if retrieving for FTC account)
	// or vip_id_alias (if retrieving for Wechat)
	stmtMember = `
	SELECT ftc_user_id AS ftc_id,
		wx_union_id AS union_id,
		member_tier AS tier,
		billing_cycle AS cycle,
		expire_date AS expire_date,
		payment_method AS payment_method,
		stripe_subscription_id AS stripe_sub_id,
		stripe_plan_id AS stripe_plan_id,
		IFNULL(auto_renewal, FALSE) AS auto_renewal,
		sub_status AS sub_status
	FROM premium.ftc_vip
	WHERE %s = ?`

	stmtLoginHistory = `
	SELECT user_id,
		auth_method AS login_method,
		client_type,
		client_version,
		INET6_NTOA(user_ip) AS login_ip,
		user_agent AS user_agent,
		created_utc AS created_at
	FROM user_db.login_history
	WHERE user_id = ?
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	stmtWxLoginHistory = `
	SELECT union_id,
		open_id,
		app_id,
		client_type,
		client_version,
		INET6_NTOA(user_ip) AS login_ip,
		user_agent AS user_agent,
		created_utc AS created_at,
		updated_utc AS updated_at
	FROM user_db.wechat_access
	WHERE union_id = ?
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	stmtOrder = `
	SELECT trade_no AS order_id,
		user_id,
		tier_to_buy AS tier,
		billing_cycle AS cycle,
	    trade_price AS price,
	    trade_amount AS amount,
		payment_method AS payment_method,
		created_utc AS created_at,
	    confirmed_utc AS confirmed_at,
		start_date AS start_date,
		end_date AS end_date,
		client_type AS client_type,
	    client_version AS client_version,
	    INET6_NTOA(user_ip_bin) AS user_ip,
		user_agent AS user_agent
	FROM premium.ftc_trade`

	stmtReaderOrders = stmtOrder + `
	WHERE user_id IN (?, ?)
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	stmtSelectOrder = stmtOrder + `
	WHERE trade_no = ?
	LIMIT 1`

	stmtGiftCard = `
	SELECT card_id AS id,
		serial_number AS serialNumber,
		DATE(FROM_UNIXTIME(expire_time)) AS expireDate,
		FROM_UNIXTIME(active_time) AS redeemedAt,
		tier AS tier,
		cycle_unit AS cycleUnit,
		cycle_value AS cycleCount
	FROM premium.scratch_card
	WHERE serial_number = ?`
)
