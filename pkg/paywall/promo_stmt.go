package paywall

// StmtCreatePromo inserts a new row into promo table.
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
    end_utc = :end_utc,
	terms_conditions = :terms_conditions`

const colsPromo = `
SELECT p.id AS promo_id,
	p.heading,
	p.cover_url,
	p.sub_heading,
	p.content,
	p.start_utc,
	p.end_utc,
	p.terms_conditions,
	p.created_utc,
	p.created_by
`

// StmtPromo retrieves a row from promo table.
const StmtPromo = colsPromo + `
FROM subs_product.paywall_promo AS p
WHERE id = ?
LIMIT 1`
