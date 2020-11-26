package subs

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/guregu/null"
)

type Charge struct {
	// The actual amount payable.
	Amount   float64 `json:"amount" db:"amount"`     // Actual price paid.
	Currency string  `json:"currency" db:"currency"` // in which currency.
}

// Order is a user's subs order
type Order struct {
	ID    string  `json:"id" db:"order_id"`
	Price float64 `json:"price" db:"price"`
	Charge
	CompoundID string      `json:"compoundId" db:"compound_id"`
	FtcID      null.String `json:"ftcId" db:"ftc_id"`
	UnionID    null.String `json:"unionId" db:"union_id"`
	PlanID     null.String `json:"planId" db:"plan_id"`
	DiscountID null.String `json:"discountId" db:"discount_id"`
	paywall.Edition
	Currency      null.String    `json:"currency"`
	CycleCount    int64          `json:"cycleCount" db:"cycle_count"`
	ExtraDays     int64          `json:"extraDays" db:"extra_days"`
	Kind          enum.OrderKind `json:"kind" db:"kind"`
	PaymentMethod enum.PayMethod `json:"payMethod" db:"payment_method"`
	TotalBalance  null.Float     `json:"totalBalance" db:"total_balance"`
	WxAppID       null.String    `json:"wxAppId" db:"wx_app_id"` // Wechat specific. Used by webhook to verify notification.
	CreatedAt     chrono.Time    `json:"createdAt" db:"created_utc"`
	ConfirmedAt   chrono.Time    `json:"confirmedAt" db:"confirmed_utc"`
	StartDate     chrono.Date    `json:"startDate" db:"start_date"`
	EndDate       chrono.Date    `json:"endDate" db:"end_date"`
}

type OrderList struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []Order `json:"data"`
	Err  error   `json:"-"`
}
