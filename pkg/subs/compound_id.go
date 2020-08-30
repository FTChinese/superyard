package subs

import "strings"

// CompoundIDs is used to get user's id from query parameter when querying orders.
type CompoundIDs struct {
	FtcID   string `query:"ftc_id"`
	UnionID string `query:"union_id"`
}

// BuildFindInSet produces a value that can be using in FIND_IN_SET(col, value).
func (c CompoundIDs) BuildFindInSet() string {
	var ids []string

	if c.FtcID != "" {
		ids = append(ids, c.FtcID)
	}

	if c.UnionID != "" {
		ids = append(ids, c.UnionID)
	}

	return strings.Join(ids, ",")
}
