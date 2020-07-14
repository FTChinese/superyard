package subs

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// MemberSnapshot saves a membership's status prior to
// placing an order.
type MemberSnapshot struct {
	ID         string              `db:"snapshot_id"`
	Reason     enum.SnapshotReason `db:"reason"`
	CreatedUTC chrono.Time         `db:"created_utc"`
	OrderID    null.String         `db:"order_id"` // Only exists when user is performing renewal or upgrading.
	Membership
}
