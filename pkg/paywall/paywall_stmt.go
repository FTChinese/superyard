package paywall

const StmtCreateBanner = `
INSERT INTO subs_product.paywall_banner
SET heading = :heading,
    cover_url = :cover_url,
    sub_heading = :sub_heading,
    content = :content,
    created_utc = UTC_TIMESTAMP(),
    updated_utc = UTC_TIMESTAMP(),
    created_by = :created_by`

const StmtBanner = `
SELECT id AS banner_id,
	heading,
	cover_url,
	sub_heading,
	content,
	created_utc,
	updated_utc
	created_by
FROM subs_product.paywall_banner
WHERE id = ?
LIMIT 1`

const StmtUpdateBanner = `
UPDATE subs_product.paywall_banner
SET heading = :heading,
    cover_url = :cover_url,
    sub_heading = :sub_heading,
    content = :content,
    updated_utc = UTC_TIMESTAMP()
WHERE id = :banner_id
LIMIT 1`

const StmtCreatePromo = `
INSERT INTO subs_product.paywall_promo
SET id = :promo_id,
	heading = :heading,
    cover_url = :cover_url,
    sub_heading = :sub_heading,
    content = :content,
    created_utc = UTC_TIMESTAMP(),
    created_by = :created_by,
    start_utc = :start_utc,
    end_utc = :end_utc`

const StmtApplyPromo = `
UPDATE subs_product.paywall_banner
SET promo_id = :promo_id,
  updated_utc = UTC_TIMESTAMP()
WHERE id = :banner_id
LIMIT 1`

const StmtDropPromo = `
UPDATE subs_product.paywall_banner
SET promo_id = NULL,
	updated_utc = UTC_TIMESTAMP()
WHERE id = ?
LIMIT 1`

const StmtPromo = `
SELECT id AS promo_id,
	heading,
	cover_url,
	sub_heading,
	content,
	start_utc,
	end_utc,
	created_utc,
	created_by
FROM subs_product.paywall_promo
WHERE id = ?
LIMIT 1`

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
