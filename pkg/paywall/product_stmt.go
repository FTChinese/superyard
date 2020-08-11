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

// StmtProduct retrieves a single product by id.
// This is used when modify an existing product.
const StmtProduct = `
SELECT id AS product_id,
    tier,
    heading,
    description,
    small_print,
    created_by,
    created_utc,
    updated_utc
FROM subs.product
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
'productId', product_id
'price', price,
'tier', tier,
'cycle', cycle,
'description', description,
'createdUtc', created_utc,
'createdBy', created_by`

const groupPlansOfProduct = `
SELECT product_id,
	JSON_ARRAYAGG(JSON_OBJECT(` + planJSONSchema + `)) AS basePlans
FROM plan
GROUP BY product_id`

const colProduct = `
SELECT prod.id,
	prod.tier,
	prod.heading,
	prod.description,
	prod.small_print,
	prod.created_by,
	prod.created_utc,
	prod.updated_utc
`

// StmtListPricedProducts retrieves a list of product.
const StmtListPricedProducts = colProduct + `
	IFNULL(plan.basePlans, JSON_ARRAY())
FROM subs.product AS prod
  	LEFT JOIN (` + groupPlansOfProduct + `) AS plan
	ON prod.id = plan.product_id
ORDER BY prod.tier ASC`
