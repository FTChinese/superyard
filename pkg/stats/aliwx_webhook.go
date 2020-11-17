package stats

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

type UnconfirmedOrder struct {
	OrderID      string         `json:"orderId" db:"order_id"`
	OrderAmount  float64        `json:"orderAmount" db:"order_amount"`
	OrderTier    enum.Tier      `json:"orderTier" db:"order_tier"`
	OrderCycle   enum.Cycle     `json:"orderCycle" db:"order_cycle"`
	Kind         enum.OrderKind `json:"kind" db:"kind"`
	CreatedUTC   chrono.Time    `json:"createdUtc" db:"created_utc"`
	ConfirmedUTC chrono.Time    `json:"confirmedUtc" db:"confirmed_utc"`
	StartDate    chrono.Date    `json:"startDate" db:"start_date"`
	EndDate      chrono.Date    `json:"endDate" db:"end_date"`
	PaymentState null.String    `json:"paymentState" db:"payment_state"`
	PaidCST      null.String    `json:"paidCst" db:"paid_cst"`
	MemberTier   enum.Tier      `json:"memberTier" db:"member_tier"`
	MemberCycle  enum.Cycle     `json:"memberCycle" db:"member_cycle"`
	ExpireDate   chrono.Date    `json:"expireDate" db:"member_expiration"`
}

type AliWxFailedList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []UnconfirmedOrder `json:"data"`
	Err  error              `json:"-"`
}
