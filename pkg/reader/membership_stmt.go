package reader

// If vip_id == vip_id_alias, the membership is purchased
// by wechat.
const membershipCols = `
SELECT vip_id AS compound_id,
	NULLIF(vip_id, vip_id_alias) AS ftc_id,
	vip_id_alias AS union_id,
	vip_type,
	expire_time,
	member_tier AS tier,
	billing_cycle AS cycle,
	expire_date,
	payment_method AS pay_method,
	ftc_plan_id,
	stripe_subscription_id AS stripe_subs_id,
	stripe_plan_id,
	auto_renewal,
	sub_status AS subs_status,
	apple_subscription_id AS apple_subs_id,
	b2b_licence_id
FROM premium.ftc_vip
`

// StmtMembership selects a reader's membership by compound id.
const StmtMembership = membershipCols + `
WHERE vip_id = ?
LIMIT 1`

// StmtMemberForOrder retrieves a reader's current membership for a particular order.
// The WHERE is used to handle a case that the order might be created when user logged in
// with wechat-only while the account is already linked to FTC account. In such a case if
// the order's user_id is wechat's union id while membership's vip_id is ftc uuid. You could never
// find out this order's current membership in this way.
const StmtMemberForOrder = membershipCols + `
WHERE ? IN (vip_id, vip_id_alias)
LIMIT 1`

const mUpsertSharedCols = `
expire_date = :expire_date,
payment_method = :pay_method,
ftc_plan_id = :ftc_plan_id,
stripe_subscription_id = :stripe_subs_id,
stripe_plan_id = :stripe_plan_id,
auto_renewal = :auto_renewal,
sub_status = :subs_status,
apple_subscription_id = :apple_subs_id,
b2b_licence_id = :b2b_licence_id`

const mUpsertCols = `
vip_type = :vip_type,
expire_time = :expire_time,
member_tier = :tier,
billing_cycle = :cycle,
` + mUpsertSharedCols

const StmtInsertMember = `
INSERT INTO premium.ftc_vip
SET vip_id = :compound_id,
	vip_id_alias = :union_id,
	ftc_user_id = :ftc_id,
	wx_union_id = :union_id,
` + mUpsertCols

const StmtUpdateMember = `
UPDATE premium.ftc_vip
SET` + mUpsertCols + `
WHERE vip_id = :compound_id
LIMIT 1`
