package paywall

// StmtCreateDiscount creates a new discount for a plan.
const StmtCreateDiscount = `
INSERT INTO subs.discount
SET id = :discount_id,
    plan_id = :plan_id,
    price_off = :price_off,
    start_utc = :start_utc,
    end_utc = :end_utc,
    created_utc = :created_utc,
    created_by = :created_by`

// StmtApplyDiscount set the plan.discount_id to the newly
// created discount.
const StmtApplyDiscount = `
UPDATE subs.plan
SET discount_id = :discount_id
WHERE id = :plan_id
LIMIT 1`

const StmtDropDiscount = `
UPDATE subs.plan
SET discount_id = NULL
WHERE id = :plan_id
LIMIT 1`
