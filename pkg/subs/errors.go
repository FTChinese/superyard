package subs

import "errors"

var (
	ErrAlreadyConfirmed   = errors.New("order already confirmed")
	ErrAlreadyUpgraded    = errors.New("already subscribed to premium edition")
	ErrValidNonAliOrWxPay = errors.New("only membership created by wxpay or alipay could be confirmed")
)
