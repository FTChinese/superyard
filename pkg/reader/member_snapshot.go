package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
)

// MemberSnapshot saves a membership's status prior to placing an order.
// Membership's current snapshot is taken prior to modifications.
// It occurs in 2 cases:
// * Manually change membership. In such case we should remember who changed it.
// * Confirm an order and membership is automatically updated.
type MemberSnapshot struct {
	ID         string              `db:"snapshot_id"`
	Reason     enum.SnapshotReason `db:"reason"`
	CreatedUTC chrono.Time         `db:"created_utc"`
	CreatedBy  null.String         `db:"created_by"`
	OrderID    null.String         `db:"order_id"` // Only exists when user is performing renewal or upgrading.
	Membership
}

func NewSnapshot(reason enum.SnapshotReason, m Membership) MemberSnapshot {
	return MemberSnapshot{
		ID:         "snp_" + rand.String(12),
		Reason:     reason,
		CreatedUTC: chrono.TimeNow(),
		CreatedBy:  null.String{},
		OrderID:    null.String{},
		Membership: m,
	}
}

// WithOrderID add the optional OrderID field.
func (s MemberSnapshot) WithOrderID(id string) MemberSnapshot {
	s.OrderID = null.StringFrom(id)

	return s
}

// WithCreator adds the optional CreatedBy field.
func (s MemberSnapshot) WithCreator(name string) MemberSnapshot {
	s.CreatedBy = null.StringFrom(name)

	return s
}

// SnapshotReasonForOrder deduces why a membership is snapshot
// when an order is confirmed and membership updated.
func SnapshotReasonForOrder(k enum.OrderKind) enum.SnapshotReason {
	switch k {
	case enum.OrderKindRenew:
		return enum.SnapshotReasonRenew

	case enum.OrderKindUpgrade:
		return enum.SnapshotReasonUpgrade

	default:
		return enum.SnapshotReasonNull
	}
}
