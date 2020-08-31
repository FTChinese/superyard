package subs

const colsOrder = `
SELECT trade_no AS order_id,
	trade_price AS price,
	trade_amount AS amount,
	user_id AS compound_id,
	ftc_user_id AS ftc_id,
	wx_union_id AS union_id,
	plan_id,
	discount_id,
	tier_to_buy AS tier,
	billing_cycle AS cycle,
	cycle_count,
	extra_days,
	category AS kind,
	payment_method,
	total_balance,
	wx_app_id,
	created_utc,
	confirmed_utc,
	start_date,
	end_date
FROM premium.ftc_trade
`

const StmtOrder = colsOrder + `
WHERE trade_no = ?
LIMIT 1`

// StmtListOrders retrieves all order belong to an FTC id, or wechat union id, or both.
const StmtListOrders = colsOrder + `
WHERE FIND_IN_SET(user_id, ?) > 0
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

const StmtConfirmOrder = `
UPDATE premium.ftc_trade
SET confirmed_utc = :confirmed_at,
	start_date = :start_date,
	end_date = :end_date
WHERE trade_no = :order_id`

const StmtProratedOrdersUsed = `
UPDATE premium.proration
SET consumed_utc = UTC_TIMESTAMP()
WHERE upgrade_order_id = ?`
