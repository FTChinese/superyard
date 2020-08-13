package paywall

const StmtCreateProduct = `
INSERT INTO subs.product
SET id = :product_id,
    tier = :tier,
    heading = :heading,
    description = :description,
    small_print = :small_print,
    created_by = :created_by,
    created_utc = UTC_TIMESTAMP(),
    updated_utc = UTC_TIMESTAMP()`

// pp refer to paywall_product table.
const colProduct = `
SELECT prod.id AS product_id,
	prod.tier,
	prod.heading,
	prod.description,
	prod.small_print,
	pp.product_id IS NOT NULL AS is_active,
	prod.created_by,
	prod.created_utc,
	prod.updated_utc
`

// StmtProduct retrieves a single product by id.
// This is used when modify an existing product,
// or showing the details of a product.
const StmtProduct = colProduct + `
FROM subs.product AS prod
	LEFT JOIN subs.paywall_product AS pp
	ON prod.id = pp.product_id
WHERE id = ?
LIMIT 1`

const StmtUpdateProduct = `
UPDATE subs.product
SET heading = :heading,
    description = :description,
    small_print = :small_print,
    updated_utc = UTC_TIMESTAMP()
WHERE id = :product_id
LIMIT 1`

const planJSONSchema = `
'id', id,
'productId', product_id,
'price', price,
'tier', tier,
'cycle', cycle,
'description', description,
'createdUtc', created_utc,
'createdBy', created_by`

const groupPlansOfProduct = `
SELECT product_id,
	JSON_ARRAYAGG(JSON_OBJECT(` + planJSONSchema + `)) AS basePlans
FROM subs.plan
GROUP BY product_id`

// StmtListPricedProducts retrieves a list of product.
const StmtListPricedProducts = colProduct + `,
	IFNULL(plan.basePlans, JSON_ARRAY()) AS plans
FROM subs.product AS prod
  	LEFT JOIN (` + groupPlansOfProduct + `) AS plan
	ON prod.id = plan.product_id
	LEFT JOIN subs.paywall_product AS pp
	ON prod.id = pp.product_id
ORDER BY prod.tier ASC`

const StmtActivateProduct = `
INSERT INTO subs.paywall_product
SET product_id = :product_id,
	tier = :tier
ON DUPLICATE KEY UPDATE
	product_id = :product_id`

const StmtHasActivePlan = `
SELECT EXISTS (
	SELECT *
	FROM subs.product_active_plans
	WHERE product_id = ?
) AS has_plan`
