package apple

const StmtListIAPSubs = `
SELECT environment,
	original_transaction_id,
	last_transaction_id,
	product_id,
	purchase_date_utc,
	expires_date_utc,
	tier,
	cycle,
	auto_renewal,
	created_utc,
	updated_utc
FROM premium.apple_subscription
ORDER BY updated_utc DESC
LIMIT ? OFFSET ?`
