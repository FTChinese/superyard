package subs

const stmtSelectOrder = `
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
FROM premium.ftc_trade
`

const StmtSelectOrder = stmtSelectOrder + `
WHERE trade_no = ?
LIMIT 1`

// StmtListOrders retrieves all order belong to an FTC id, or wechat union id, or both.
const StmtListOrders = stmtSelectOrder + `
WHERE FIND_IN_SET(user_id, ?) > 0
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

const StmtConfirmOrder = `
UPDATE premium.ftc_trade
SET confirmed_utc = :confirmed_at,
	start_date = :start_date,
	end_date = :end_date
WHERE trade_no = :order_id`
