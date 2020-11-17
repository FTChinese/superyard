package subs

import (
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
