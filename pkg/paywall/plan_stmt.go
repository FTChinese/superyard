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
	
	p.created_utc,
	p.created_by`

// StmtPlan selects a single plan.
const StmtPlan = colPlan + `
FROM subs_product.plan AS p
WHERE p.id = ?
LIMIT 1`

const StmtActivatePlan = `
INSERT INTO subs_product.product_active_plans
SET plan_id = :plan_id,
	product_id = :product_id,
	cycle = :cycle
ON DUPLICATE KEY UPDATE
	plan_id = :plan_id`

const colDiscountedPlan = colPlan + `,
d.id AS discount_id,
d.price_off,
d.percent,
d.start_utc,
d.end_utc
`

// StmtPlansOfProduct selects all plans under a product.
const StmtPlansOfProduct = colDiscountedPlan + `,
	a.plan_id IS NOT NULL AS is_active
FROM subs_product.plan AS p
    LEFT JOIN subs_product.discount AS d
	ON p.discount_id = d.id
	LEFT JOIN subs_product.product_active_plans AS a
	ON p.id = a.plan_id
WHERE p.product_id = ?
ORDER BY p.cycle DESC`
