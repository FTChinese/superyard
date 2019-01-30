package user

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// Order is a user's subscription order
type Order struct {
	ID            string         `json:"orderId"`
	Tier          enum.Tier      `json:"tier"`
	Cycle         enum.Cycle     `json:"cycle"`
	ListPrice     float64        `json:"listPrice"`
	NetPrice      float64        `json:"netPrice"`
	PaymentMethod enum.PayMethod `json:"payMethod"`
	CreatedAt     chrono.Time `json:"createdAt"`
	ConfirmedAt   chrono.Time `json:"confirmedAt"`
	StartDate     chrono.Date `json:"startDate"`
	EndDate       chrono.Date `json:"endDate"`
	ClientType    null.String      `json:"clientType"`
	ClientVersion null.String      `json:"clientVersion"`
	UserIP        null.String      `json:"userIp"`
}
