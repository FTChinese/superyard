package stats

const fromAliUnconfirmed = `
FROM premium.log_ali_notification AS a
    LEFT JOIN premium.ftc_trade AS o
    ON a.ftc_order_id = o.trade_no
    LEFT JOIN premium.ftc_vip AS m
    ON o.user_id = m.vip_id
WHERE o.trade_no IS NOT NULL
    AND o.confirmed_utc IS NULL
    AND a.trade_status = 'TRADE_SUCCESS'
`

const fromWxUnconfirmed = `
FROM premium.log_wx_notification AS w
    LEFT JOIN premium.ftc_trade AS o
    ON w.ftc_order_id = o.trade_no
    LEFT JOIN premium.ftc_vip AS m
    ON o.user_id = m.vip_id
WHERE o.trade_no IS NOT NULL
    AND o.confirmed_utc IS NULL
    AND w.result_code = 'SUCCESS'
`

const StmtCountAliUnconfirmed = `
SELECT COUNT(*)
` + fromAliUnconfirmed

const StmtAliUnconfirmed = `
SELECT o.trade_no AS order_id,
	o.trade_amount AS order_amount,
    o.tier_to_buy AS order_tier,
    o.billing_cycle AS order_cycle,
    o.category AS kind,
    o.created_utc AS created_utc,
    o.confirmed_utc AS confirmed_utc,
    o.start_date AS start_date,
    o.end_date AS end_date,

    a.trade_status AS payment_state,
    a.paid_cst AS paid_cst,
    
    m.member_tier AS member_tier,
    m.billing_cycle AS member_cycle,
    m.expire_date AS member_expiration
` + fromAliUnconfirmed + `
ORDER BY o.created_utc DESC
LIMIT ? OFFSET ?`

const StmtCountWxUnconfirmed = `
SELECT COUNT(*)
` + fromWxUnconfirmed

const StmtWxUnconfirmed = `
SELECT o.trade_no AS order_id,
    o.trade_amount AS order_amount,
    o.tier_to_buy AS order_tier,
    o.billing_cycle AS order_cycle,
    o.category AS kind,
    o.created_utc,
    o.confirmed_utc,
    o.start_date,
    o.end_date,

    w.result_code AS payment_state,
    w.time_end AS paid_cst,
    
    m.member_tier AS member_tier,
    m.billing_cycle AS member_cycle,
    m.expire_date AS member_expiration
` + fromWxUnconfirmed + `
ORDER BY o.created_utc DESC
LIMIT ? OFFSET ?`
