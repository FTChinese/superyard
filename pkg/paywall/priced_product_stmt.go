package paywall

// StmtListPricedProducts retrieves a list of product, and each product has a list of plans.
const StmtListPricedProducts = colProduct + `,
	IFNULL(groupedPlans.plan_count, 0) AS plan_count
FROM subs_product.product AS prod
  	LEFT JOIN (
		SELECT plan.product_id, 
			COUNT(*) AS plan_count
		FROM subs_product.plan AS plan
		GROUP BY plan.product_id
	) AS groupedPlans
	ON prod.id = groupedPlans.product_id
	LEFT JOIN subs_product.paywall_product AS pp
	ON prod.id = pp.product_id
ORDER BY prod.tier ASC`
