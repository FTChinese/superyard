package paywall

const StmtCreatePlan = `
INSERT INTO subs_product.plan
SET id = :plan_id,
    product_id = :product_id,
    price = :price,
    tier = :tier,
    cycle = :cycle,
    description = :description,
    created_utc = UTC_TIMESTAMP(),
    created_by = :created_by`

const colPlan = `
SELECT p.id AS plan_id,
	p.product_id,
	p.price,
	p.tier,
	p.cycle,
	p.description,
	a.plan_id IS NOT NULL AS is_active,
	p.created_utc,
	p.created_by`

// StmtListPlansInUse retrieves all plans that is used on paywall.
// This is used by client UI to show a list of plans
// when creating/updating a membership for wx or alipay.
const StmtListPlansOnPaywall = colPlan + `
FROM subs_product.product_active_plans AS a
	LEFT JOIN subs_product.paywall_product AS pp
		ON a.product_id = pp.product_id
	LEFT JOIN subs_product.plan AS p
		ON a.plan_id = p.id
WHERE pp.product_id IS NOT NULL
	AND p.id IS NOT NULL
`

const StmtPaywallPlan = StmtListPlansOnPaywall + `
AND pp.tier = ?
AND a.cycle = ?
LIMIT 1`

const StmtActivatePlan = `
INSERT INTO subs_product.product_active_plans
SET plan_id = :plan_id,
	product_id = :product_id,
	cycle = :cycle
ON DUPLICATE KEY UPDATE
	plan_id = :plan_id`
