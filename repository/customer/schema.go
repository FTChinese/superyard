package customer

import "fmt"

const (
	selectMember = `
	SELECT id AS member_id, 
		vip_id AS compound_id,
		NULLIF(vip_id, vip_id_alias) AS ftc_id,
		vip_id_alias AS union_id,
		vip_type,
		expire_time,
		member_tier AS tier,
		billing_cycle AS cycle,
		expire_date,
		payment_method,
		stripe_subscription_id AS stripe_sub_id,
		auto_renewal,
		sub_status
	FROM premium.ftc_vip`

	selectMemberByID = selectMember + `
	WHERE id = ?
	LIMIT 1`

	memberForEmail = selectMember + `
	WHERE vip_id = ?
	LIMIT 1`

	memberForWx = selectMember + `
	WHERE vip_id_alias = ?
	LIMIT 1`

	stmtUpdateMember = `
	UPDATE premium.ftc_vip
	SET vip_type = :vip_type,
		expire_time = :expire_time,
		member_tier = :tier,
		billing_cycle = :cycle,
		expire_date = :expire_date,
		payment_method = :payment_method,
		stripe_subscription_id = :stripe_sub_id,
		stripe_plan_id = :stripe_plan_id,
		auto_renewal = :auto_renewal,
		sub_status = :sub_status
	WHERE id = :member_id
	LIMIT 1`

	stmtInsertMember = `
	INSERT INTO premium.ftc_vip
	SET id = :member_id,
		vip_id = :compound_id,
		vip_id_alias = :union_id,
		vip_type = :vip_type,
		expire_time = :expire_time,
		ftc_user_id = :ftc_id,
		wx_union_id = :union_id,
		member_tier = :tier,
		billing_cycle = :cycle,
		expire_date = :expire_date,
		payment_method = :payment_method,
		stripe_subscription_id = :stripe_sub_id,
		stripe_plan_id = :stripe_plan_id,
		auto_renewal = :auto_renewal,
		sub_status = :sub_status`

	stmtDeleteMember = `
	DELETE FROM premium.ftc_vip
	WHERE id = :member_id
	LIMIT 1`

	insertMemberSnapshot = `
	INSERT INTO %s.member_snapshot
	SET id = :snapshot_id,
		reason = :reason,
		created_utc = UTC_TIMESTAMP(),
		member_id = :member_id,
		compound_id = :compound_id,
		ftc_user_id = :ftc_id,
		wx_union_id = :union_id,
		expire_date = :expire_date,
		payment_method = :payment_method,
		stripe_subscription_id = :stripe_sub_id,
		stripe_plan_id = :stripe_plan_id,
		auto_renewal = :auto_renewal,
		sub_status = :sub_status`

	// Orders
	// ---------------------------
	stmtSelectOrder = `
	SELECT trade_no AS order_id,
		user_id AS compound_id,
		ftc_user_id AS ftc_id,
		wx_union_id AS union_id,
		trade_price AS price,
		trade_amount AS amount,
		tier_to_buy AS tier,
		billing_cycle AS cycle,
		cycle_count,
		extra_days,
		category AS usage_type,
		payment_method,
		created_utc AS created_at,
		confirmed_utc AS confirmed_at,
		start_date,
		end_date,
		upgrade_id,
		member_snapshot_id
	FROM premium.ftc_trade`

	stmtAnOrder = stmtSelectOrder + `
	WHERE trade_no = ?
	LIMIT 1`

	stmtInsertOrder = `
	INSERT INTO ftc_trade.ftc_trade
	SET trade_no = :order_id,
		user_id = :compound_id,
		ftc_user_id = :ftc_id,
		wx_union_id = :union_id,
		trade_price = :price,
		trade_amount = :amount,
		tier_to_buy = :tier,
		billing_cycle = :cycle,
		cycle_count = :cycle_count,
		extra_days = :extra_days,
		category = :usage_type,
		payment_method = :payment_method,
		wx_app_id = wx_app_id,
		created_utc = UTC_TIMESTAMP(),
		upgrade_id = :upgrade_id,
		member_snapshot_id = :member_snapshot_id`

	stmtConfirmOrder = `
	UPDATE premium.ftc_trade
	SET confirmed_utc = :confirmed_at,
		start_date = :start_date,
		end_date = :end_date
	WHERE trade_no = :id`

	// Retrieve the client when user creates an order.
	stmtOrderClient = `
	SELECT client_type,
	    client_version,
	    INET6_NTOA(user_ip_bin) AS user_ip,
		user_agent
	FROM premium.client
	WHERE order_id = ?
	LIMIT 1`

	// ---------------------
	// Login history
	stmtLoginHistory = `
	SELECT user_id,
		auth_method AS login_method,
		client_type,
		client_version,
		INET6_NTOA(user_ip) AS user_ip,
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
		INET6_NTOA(user_ip) AS user_ip,
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

func stmtListOrders(bindVar string) string {
	return fmt.Sprintf(`
	%s
	WHERE user_id IN (%s)
	ORDER BY created_utc DESC
	LIMIT ? OFFSET ?`, stmtSelectOrder, bindVar)
}
