package paywall

const StmtInsertPlan = `
INSERT INTO subs.plan
SET id = :plan_id,
    product_id = :product_id,
    price = :price,
    tier = :tier,
    cycle = :cycle,
    description = :description,
    created_utc = UTC_TIMESTAMP(),
    created_by = :created_by`

const colPlan = `
SELECT p.id,
	p.product_id,
	p.price,
	p.tier,
	p.cycle,
	p.description,
	p.created_utc,
	p.created_by,
	d.id AS discount_id,
	d.price_off,
	d.percent,
	d.start_utc,
	d.end_utc
`

// StmtSelectPlan selects a single plan.
const StmtSinglePlan = colPlan + `
FROM subs.plan AS p
    LEFT JOIN subs.discount AS d
	ON p.discount_id = d.id
WHERE p.id = ?
LIMIT 1`

// StmtPlansOfProduct selects all plans under a product.
const StmtPlansOfProduct = colPlan + `
	a.plan_id IS NOT NULL AS is_active
FROM subs.plan AS p
    LEFT JOIN subs.discount AS d
	ON p.discount_id = d.id
	LEFT JOIN subs.product_active_plans AS a
	ON p.plan_id = pap.plan_id
WHERE p.product_id = ?
ORDER BY p.cycle DESC`

const StmtActivatePlan = `
INSERT INTO subs.product_active_plans
SET plan_id = :plan_id
  product_id = :product_id
  cycle = :cycle
ON DUPLICATE KEY UPDATE
  plan_id = :plan_id`
