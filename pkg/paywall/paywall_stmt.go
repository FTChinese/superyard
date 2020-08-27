package paywall

// StmtPaywallPromo retrieves a promo that is used on paywall.
const StmtPaywallPromo = colsPromo + `
FROM subs_product.paywall_banner AS b
	LEFT JOIN subs_product.paywall_promo AS p
	ON b.promo_id = p.id
WHERE b.id = ?
LIMIT 1`

// StmtPaywallProducts retrieves the products shown on paywall.
const StmtPaywallProducts = colProduct + `
FROM subs_product.paywall_product AS pp
	LEFT JOIN subs_product.product AS prod
	ON pp.product_id = prod.id
WHERE prod.id IS NOT NULL
ORDER BY prod.tier ASC`

// StmtPaywallPlans selects all active plans of products which are listed on paywall.
// The plans has discount attached.
const StmtPaywallPlans = colDiscountedPlan + `
FROM subs_product.product_active_plans AS a
	LEFT JOIN subs_product.plan AS p
	ON a.plan_id = p.id
	LEFT JOIN subs_product.discount AS d
	ON p.discount_id = d.id
	LEFT JOIN subs_product.paywall_product AS pp
	ON a.product_id = pp.product_id
WHERE p.id IS NOT NULL
	AND pp.product_id IS NOT NULL
ORDER BY cycle DESC`
