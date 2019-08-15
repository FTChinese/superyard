package customer

type MemberColumn string

const (
	MemberColumnID    MemberColumn = "id"
	MemberColumnFtcID              = "vip_id"
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

	// Subscription
	// -------------------------
	// Select membership by either vip_id
	// (if retrieving for FTC account)
	// or vip_id_alias (if retrieving for Wechat)
	stmtMember = `
	SELECT id,
		ftc_user_id AS ftc_id,
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

	stmtInsertMember = `
	INSERT INTO premium.ftc_vip
	SET id = :id,
		vip_id = :compound_id,
		vip_id_alias = :union_id,
		ftc_user_id = :ftc_id,
		wx_union_id = :union_id,
		vip_type = :vip_type,
		expire_time = :expire_time,
		member_tier = :tier,
		billing_cycle = :cycle,
		expire_date = :expire_date,
		payment_method = :payment_method,
		stripe_subscription_id = :stripe_sub_id,
		stripe_plan_id = :stripe_plan_id,
		auto_renewal = :auto_renewal,
		sub_status = :sub_status`

	stmtUpdateMember = `
	UPDATE premium.ftc_vip
	SET vip_type = :vip_type,
		expire_time = :expire_time,
		member_tier = :tier,
		billing_cycle = :cycle,
		expire_date = :expire_date,
		payment_method = :payment_method
		stripe_subscription_id = :stripe_sub_id,
		stripe_plan_id = :stripe_plan_id,
		auto_renewal = :auto_renewal,
		sub_status = :sub_status
	WHERE id = ?
	LIMIT 1`

	stmtDeleteMember = `
	`

	// Orders
	// ---------------------------
	stmtSelectOrder = `
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

	stmtListOrders = stmtSelectOrder + `
	WHERE user_id IN (?, ?)
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`

	stmtSelectOneOrder = stmtSelectOrder + `
	WHERE trade_no = ?
	LIMIT 1`

	stmtCreateOrder = `
	INSERT INTO premium.ftc_trade
	SET trade_no = :id,
		trade_price = :price,
		trade_amount = :amount,
		user_id = :compound_id,
		ftc_user_id = :ftc_id,
		wx_union_id = :union_id,
		tier_to_buy = :tier,
		billing_cycle = :cycle,
		cycle_count = :cycle_count,
		extra_days = :extra_days,
		category = :usage,
		last_upgrade_id = :last_upgrade_id,
		payment_method =:payment_method,
		created_utc = UTC_TIMESTAMP(),
		confirmed_utc = UTC_TIMESTAMP(),
		start_date = :start_date,
		end_date = :end_date`

	stmtConfirmOrder = `
	UPDATE premium.ftc_trade
	SET confirmed_utc = UTC_TIMESTAMP()
		start_date = :start_date,
		end_date = :end_date
	WHERE trade_no = :id`

	// ---------------------
	// Login history
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
