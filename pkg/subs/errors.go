package subs

import "errors"

var (
	ErrAlreadyConfirmed = errors.New("order already confirmed")
	ErrAlreadyUpgraded  = errors.New("already subscribed to premium edition")
)
