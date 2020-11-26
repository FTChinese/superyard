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
`

const StmtOrder = colsOrder + `
FROM premium.ftc_trade
WHERE trade_no = ?
LIMIT 1`

const fromListOrder = `
FROM premium.ftc_trade
WHERE FIND_IN_SET(user_id, ?) > 0
`

// StmtListOrders retrieves all order belong to an FTC id, or wechat union id, or both.
const StmtListOrders = colsOrder + fromListOrder + `
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

const StmtCountOrder = `
SELECT COUNT(*) AS row_count
` + fromListOrder
