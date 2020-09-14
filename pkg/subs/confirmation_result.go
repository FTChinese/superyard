package subs

import "github.com/FTChinese/superyard/pkg/reader"

// ConfirmationResult contains the result of confirming an order.
// The value of Order and Membership are inter-dependent on each other:
// You have to use a Membership's existing expiration date
// to determine this order's subscribed period;
// then you have to use the confirmed order's period
// to update the Membership to a new expiration date.
// You also have to keep a snapshot of existing membership
// for backup purpose.
type ConfirmationResult struct {
	Order      Order                 // The confirmed order.
	Membership reader.Membership     // The updated membership. Might be zero value.
	Snapshot   reader.MemberSnapshot // Snapshot of previous membership. Might be empty is Membership is empty.
}
