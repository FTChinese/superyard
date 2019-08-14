package reader

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

// Order is a user's subs order
type Order struct {
	ID            string         `json:"orderId" db:"order_id"`
	UserID        string         `json:"userId" db:"user_id"`
	Tier          enum.Tier      `json:"tier" db:"tier"`
	Cycle         enum.Cycle     `json:"cycle" db:"cycle"`
	ListPrice     float64        `json:"price" db:"price"`
	NetPrice      float64        `json:"amount" db:"amount"`
	PaymentMethod enum.PayMethod `json:"payMethod" db:"payment_method"`
	CreatedAt     chrono.Time    `json:"createdAt" db:"created_at"`
	ConfirmedAt   chrono.Time    `json:"confirmedAt" db:"confirmed_at"`
	StartDate     chrono.Date    `json:"startDate" db:"start_date"`
	EndDate       chrono.Date    `json:"endDate" db:"end_date"`
	ClientType    enum.Platform  `json:"clientType" db:"client_type"`
	ClientVersion null.String    `json:"clientVersion" db:"client_version"`
	UserIP        null.String    `json:"userIp" db:"user_ip"`
	UserAgent     null.String    `json:"userAgent" db:"user_agent"`
}
