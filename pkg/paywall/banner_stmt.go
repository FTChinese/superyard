package paywall

// StmtCreateBanner insert a row in banner table.
const StmtCreateBanner = `
INSERT INTO subs_product.paywall_banner
SET heading = :heading,
    cover_url = :cover_url,
    sub_heading = :sub_heading,
    content = :content,
    created_utc = UTC_TIMESTAMP(),
    updated_utc = UTC_TIMESTAMP(),
    created_by = :created_by`

// StmtBanner retrieves a banner by id. The id is always 1.
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

// StmtUpdateBanner updates the existing banner.
const StmtUpdateBanner = `
UPDATE subs_product.paywall_banner
SET heading = :heading,
    cover_url = :cover_url,
    sub_heading = :sub_heading,
    content = :content,
    updated_utc = UTC_TIMESTAMP()
WHERE id = :banner_id
LIMIT 1`

// StmtApplyPromo sets a promo's id to banner.
const StmtApplyPromo = `
UPDATE subs_product.paywall_banner
SET promo_id = :promo_id,
  updated_utc = UTC_TIMESTAMP()
WHERE id = :banner_id
LIMIT 1`

// StmtDropPromo removes the promo id from banner.
const StmtDropPromo = `
UPDATE subs_product.paywall_banner
SET promo_id = NULL,
	updated_utc = UTC_TIMESTAMP()
WHERE id = ?
LIMIT 1`
