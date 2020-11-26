package reader

const InsertMemberSnapshot = `
INSERT INTO premium.member_snapshot
SET id = :snapshot_id,
	created_by = :created_by,
	created_utc = UTC_TIMESTAMP(),
	order_id = :order_id,
	compound_id = :compound_id,
	ftc_user_id = :ftc_id,
	wx_union_id = :union_id,
	tier = :tier,
	cycle = :cycle,
` + mUpsertSharedCols

const fromSnapshot = `
FROM premium.member_snapshot
WHERE FIND_IN_SET(compound_id, ?) > 0
`

const StmtMemberSnapshots = `
SELECT id AS snapshot_id,
	created_utc,
	created_by
	order_id,
	compound_id,
	ftc_user_id AS ftc_id,
	wx_union_id AS union_id,
	tier,
	cycle,
` + colMemberShared + fromSnapshot + `
ORDER BY created_utc DESC
LIMIT ? OFFSET ?`

const StmtCountMemberSnapshot = `
SELECT COUNT(*) AS row_count
` + fromSnapshot
