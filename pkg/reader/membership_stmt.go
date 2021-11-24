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

// StmtSelectMember selects a reader's membership by compound id.
const StmtSelectMember = colMembership + `
WHERE ? IN (vip_id, vip_id_alias)
LIMIT 1`

const StmtIAPMember = colMembership + `
WHERE apple_subscription_id = ?
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
b2b_licence_id = :b2b_licence_id
`

const mUpsertCols = `
vip_type = :vip_type,
expire_time = :expire_time,
member_tier = :tier,
billing_cycle = :cycle,
` + mUpsertSharedCols

const StmtCreateMember = `
INSERT INTO premium.ftc_vip
SET vip_id = :compound_id,
	vip_id_alias = :union_id,
	ftc_user_id = :ftc_id,
	wx_union_id = :union_id,
` + mUpsertCols

// StmtDeleteMember deletes the membership under a sandbox account.
// Never delete a real user's membership.
const StmtDeleteMember = `
DELETE FROM premium.ftc_vip
WHERE vip_id = ?
LIMIT 1`
