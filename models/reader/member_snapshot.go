package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
)

// MemberSnapshot saves a membership's status prior to
// placing an order.
type MemberSnapshot struct {
	ID         string         `db:"snapshot_id"`
	Reason     SnapshotReason `db:"reason"`
	CreatedUTC chrono.Time    `db:"created_utc"`
	Membership
}

func NewMemberSnapshot(m Membership, r SnapshotReason) MemberSnapshot {
	return MemberSnapshot{
		ID:         "snp_" + rand.String(12),
		Reason:     r,
		CreatedUTC: chrono.TimeNow(),
		Membership: m,
	}
}
