package paywall

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
FROM subs_product.plan
GROUP BY product_id`

// StmtListPricedProducts retrieves a list of product, and each product has a list of plans.
const StmtListPricedProducts = colProduct + `,
	IFNULL(plan.basePlans, JSON_ARRAY()) AS plans
FROM subs_product.product AS prod
  	LEFT JOIN (` + groupPlansOfProduct + `) AS plan
	ON prod.id = plan.product_id
	LEFT JOIN subs_product.paywall_product AS pp
	ON prod.id = pp.product_id
ORDER BY prod.tier ASC`
