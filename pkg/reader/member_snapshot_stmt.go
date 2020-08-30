package reader

const InsertMemberSnapshot = `
INSERT INTO premium.member_snapshot
SET id = :snapshot_id,
	reason = :reason,
	created_utc = UTC_TIMESTAMP(),
	created_by = :created_by,
	order_id = :order_id,
	compound_id = :compound_id,
	ftc_user_id = :ftc_id,
	wx_union_id = :union_id,
	tier = :tier,
	cycle = :cycle,
` + mUpsertSharedCols

const StmtMemberSnapshot = `
SELECT id,
	reason,
	created_utc,
	created_by
	order_id,
	compound_id,
	ftc_user_id AS ftc_id,
	wx_union_id AS union_id,
	tier,
	cycle,
	expire_date,
	payment_method,
	ftc_plan_id,
	stripe_subs_id,
	stripe_plan_id,
	auto_renewal,
	sub_status AS subs_status,
	apple_subscription_id AS apple_subs_id,
	b2b_licence_id
FROM premium.member_snapshot
WHERE FIND_IN_SET(compound_id, ?)`
