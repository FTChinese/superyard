package subs

import (
	"github.com/FTChinese/go-rest/chrono"
	"github.com/FTChinese/go-rest/enum"
	"github.com/guregu/null"
)

const StmtAliPayload = `
SELECT trade_number,
	ftc_order_id,
	trade_status,
	total_amount,
	receipt_amount,
	notified_cst,
	created_cst,
	paid_cst,
	closed_cst
FROM premium.log_ali_notification
WHERE ftc_order_id = ?
ORDER BY notified_cst DESC`

type AliPayload struct {
	TransactionID string      `json:"transactionId" db:"trade_number"`
	FtcOrderID    string      `json:"ftcOrderId" db:"ftc_order_id"`
	TradeStatus   string      `json:"tradeStatus" db:"trade_status"`
	TotalAmount   string      `json:"totalAmount" db:"total_amount"`
	ReceiptAmount null.String `json:"receiptAmount" db:"receipt_amount"`
	NotifiedCST   string      `json:"notifiedCst" db:"notified_cst"`
	CreatedCST    string      `json:"createdCst" db:"created_cst"`
	PaidCST       string      `json:"paidCst" db:"paid_cst"`
	ClosedCST     null.String `json:"closedCst" db:"closed_cst"`
}

const StmtWxPayload = `
SELECT return_code,
	return_message,
	result_code,
	error_code,
	error_description,
	transaction_id,
	ftc_order_id,
	trade_type,
	total_fee,
	time_end
FROM premium.log_wx_notification
WHERE ftc_order_id = ?
ORDER BY time_end DESC`

type WxPayload struct {
	ReturnCode    string      `json:"returnCode" db:"return_code"`
	ReturnMessage null.String `json:"returnMessage" db:"return_message"`
	ResultCode    string      `json:"resultCode" db:"result_code"`
	ErrorCode     null.String `json:"errorCode" db:"error_code"`
	ErrorDesc     null.String `json:"errorDesc" db:"error_description"`
	TransactionID string      `json:"transactionId" db:"transaction_id"`
	FtcOrderID    string      `json:"ftcOrderId" db:"ftc_order_id"`
	TradeType     string      `json:"tradeType" db:"trade_type"`
	TotalAmount   int64       `json:"totalAmount" db:"total_fee"`
	PaidCST       string      `json:"paidCst" db:"time_end"`
}

type UnconfirmedOrder struct {
	OrderID      string         `json:"orderId" db:"order_id"`
	OrderTier    enum.Tier      `json:"orderTier" db:"order_tier"`
	OrderCycle   enum.Cycle     `json:"orderCycle" db:"order_cycle"`
	Kind         enum.OrderKind `json:"kind" db:"kind"`
	CreatedUTC   chrono.Time    `json:"createdUtc" db:"created_utc"`
	ConfirmedUTC chrono.Time    `json:"confirmedUtc" db:"confirmed_utc"`
	StartDate    chrono.Date    `json:"startDate" db:"start_date"`
	EndDate      chrono.Date    `json:"endDate" db:"end_date"`
	PaymentState null.String    `json:"paymentState" db:"payment_state"`
	PaidAmount   null.String    `json:"paidAmount" db:"paid_amount"`
	PaidCST      null.String    `json:"paidCst" db:"paid_cst"`
	MemberTier   enum.Tier      `json:"memberTier" db:"member_tier"`
	MemberCycle  enum.Cycle     `json:"memberCycle" db:"member_cycle"`
	ExpireDate   chrono.Date    `json:"expireDate" db:"member_expiration"`
}

const StmtAliUnconfirmed = `
SELECT
    o.trade_no AS order_id,
    o.tier_to_buy AS order_tier,
    o.billing_cycle AS order_cycle,
    o.category AS kind,
    o.created_utc AS created_utc,
    o.confirmed_utc AS confirmed_utc,
    o.start_date AS start_date,
    o.end_date AS end_date,

    a.trade_status AS payment_state,
    a.receipt_amount AS paid_amount,
    a.paid_cst AS paid_cst,
    
    m.member_tier AS member_tier,
    m.billing_cycle AS member_cycle,
    m.expire_date AS member_expiration
FROM premium.log_ali_notification AS a
    LEFT JOIN premium.ftc_trade AS o
    ON a.ftc_order_id = o.trade_no
    LEFT JOIN premium.ftc_vip AS m
    ON o.user_id = m.vip_id
WHERE o.trade_no IS NOT NULL
    AND o.confirmed_utc IS NULL
    AND a.trade_status = 'TRADE_SUCCESS'
ORDER BY o.created_utc DESC`

const StmtWxUnconfirmed = `
SELECT 
    o.trade_no AS order_id,
    o.user_id AS compound_id,
    o.ftc_user_id,
    o.wx_union_id,
    o.tier_to_buy AS order_tier,
    o.billing_cycle AS order_cycle,
    o.category AS kind,
    o.created_utc,
    o.confirmed_utc,
    o.start_date,
    o.end_date,

    w.result_code AS payment_state,
    w.total_fee AS paid_amount,
    w.time_end AS paid_cst,
    
    m.member_tier AS member_tier,
    m.billing_cycle AS member_cycle,
    m.expire_date AS member_expiration
FROM premium.log_wx_notification AS w
    LEFT JOIN premium.ftc_trade AS o
    ON w.ftc_order_id = o.trade_no
    LEFT JOIN premium.ftc_vip AS m
    ON o.user_id = m.vip_id
WHERE o.trade_no IS NOT NULL
    AND o.confirmed_utc IS NULL
    AND w.result_code = 'SUCCESS'
ORDER BY o.created_utc DESC`
