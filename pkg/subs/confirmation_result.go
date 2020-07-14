package subs

// ConfirmationResult contains the result of confirming an order.
// The value of Order and Membership are inter-dependent on each other:
// You have to use a Membership's existing expiration date
// to determine this order's subscribed period;
// then you have to use the confirmed order's period
// to update the Membership to a new expiration date.
// You also have to keep a snapshot of existing membership
// for backup purpose.
type ConfirmationResult struct {
	Order      Order          // The confirmed order.
	Membership Membership     // The updated membership.
	Snapshot   MemberSnapshot // // Snapshot of previous membership
}
