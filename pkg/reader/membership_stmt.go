package reader

const colMemberShared = `
expire_date,
payment_method AS pay_method,
ftc_plan_id,
stripe_subscription_id AS stripe_subs_id,
stripe_plan_id,
IFNULL(auto_renewal, FALSE) AS auto_renewal,
sub_status AS subs_status,
apple_subscription_id AS apple_subs_id,
b2b_licence_id
`

// If vip_id == vip_id_alias, the membership is purchased
// by wechat.
const colMembership = `
SELECT vip_id AS compound_id,
	NULLIF(vip_id, vip_id_alias) AS ftc_id,
	vip_id_alias AS union_id,
	vip_type,
	expire_time,
	member_tier AS tier,
	billing_cycle AS cycle,
` + colMemberShared + `
FROM premium.ftc_vip
`

const StmtIAPMember = colMembership + `
WHERE apple_subscription_id = ?
LIMIT 1`
