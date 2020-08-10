package paywall

const StmtCreateBanner = `
INSERT INTO subs.paywall_banner
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
	create_utc,
	updated_utc
	created_by
FROM subs.paywall_banner
WHERE id = 1
LIMIT 1`

const StmtUpdateBanner = `
UPDATE subs.paywall_banner
SET heading = :heading,
    cover_url = :cover_url,
    sub_heading = :sub_heading,
    content = :content,
    updated_utc = UTC_TIMESTAMP(),
WHERE id = :banner_id
LIMIT 1`

const StmtCreatePromo = `
INSERT INTO subs.paywall_promo
SET heading = :heading,
    cover_url = :cover_url,
    sub_heading = :sub_heading,
    content = :content,
    created_utc = UTC_TIMESTAMP(),
    created_by = :created_by
    promo_start_utc = :start_utc
    promo_end_utc = :end_utc`

const StmtApplyPromo = `
UPDATE subs.paywall_banner
SET promo_id = :promo_id,
  updated_utc = UTC_TIMESTAMP()
WHERE id = 1
LIMIT 1`

const StmtPromo = `
SELECT id AS promo_id,
	heading,
	cover_url,
	sub_heading,
	content,
	start_utc,
	end_utc,
	create_utc,
	created_by
FROM subs.paywall_banner
WHERE id = ?
LIMIT 1`

const StmtActiveProducts = colProduct + `
FROM subs.product AS prod
ORDER BY prod.tier ASC`

const StmtActivePlans = colPlan + `
FROM subs.product_active_plans AS a
	LEFT JOIN subs.plan AS p
	ON a.plan_id = p.plan_id
	LEFT JOIN subs.discount AS d
	ON p.discount_id = d.id
WHERE FIND_IN_SET(a.product_id, ?) > 0
	AND p.id IS NOT NULL
ORDER BY cycle DESC`
