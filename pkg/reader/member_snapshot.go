package reader

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/rand"
	"github.com/guregu/null"
)

// MemberSnapshot saves a membership's status prior to placing an order.
// Membership's current snapshot is taken prior to modifications.
// It occurs in 2 cases:
// * Manually change membership. In such case we should remember who changed it.
// * Confirm an order and membership is automatically updated.
type MemberSnapshot struct {
	ID         string      `json:"id" db:"snapshot_id"`
	CreatedBy  null.String `json:"createdBy" db:"created_by"`
	CreatedUTC chrono.Time `json:"createdUtc" db:"created_utc"`
	OrderID    null.String `json:"orderId" db:"order_id"` // Only exists when user is performing renewal or upgrading.
	Membership
}

func NewSnapshot(m Membership) MemberSnapshot {
	return MemberSnapshot{
		ID:         "snp_" + rand.String(12),
		CreatedBy:  null.String{},
		CreatedUTC: chrono.TimeNow(),
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

type MemberRevisions struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []MemberSnapshot `json:"data"`
	Err  error            `json:"-"`
}
