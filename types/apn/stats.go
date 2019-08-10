package apn

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/guregu/null"
)

type TimeZone struct {
	ZoneName    string      `json:"name"`
	Offset      null.String `json:"offset"`
	DeviceCount null.Int    `json:"deviceCount"`
}

type Device struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type InvalidDevice struct {
	Reason string `json:"reason"`
	Name   string `json:"name"`
	Count  int64  `json:"count"`
}

type TestDevice struct {
	ID          string      `json:"id"`
	Token       string      `json:"token"`
	Description null.String `json:"description"`
	OwnedBy     null.String `json:"ownedBy"`
	CreatedAt   chrono.Time `json:"createdAt"`
}
