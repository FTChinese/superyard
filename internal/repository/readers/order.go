package readers

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/subs"
)

// ListOrders retrieves a user's orders.
// Turn reader's possible into a format used in
// MySQL function FIND_IN_SET.
func (env Env) ListOrders(ids subs.CompoundIDs, p gorest.Pagination) ([]subs.Order, error) {

	var orders = make([]subs.Order, 0)

	err := env.db.Select(
		&orders,
		subs.StmtListOrders,
		ids.BuildFindInSet(),
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// RetrieveOrder retrieves a single order by trade_no column.
func (env Env) RetrieveOrder(orderID string) (subs.Order, error) {
	var order subs.Order

	err := env.db.Get(&order, subs.StmtOrder, orderID)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (env Env) AliWebhook(orderID string) ([]subs.AliPayload, error) {
	var p = make([]subs.AliPayload, 0)

	err := env.db.Select(&p, subs.StmtAliPayload, orderID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (env Env) WxWebhook(orderID string) ([]subs.WxPayload, error) {
	var p = make([]subs.WxPayload, 0)

	err := env.db.Select(&p, subs.StmtWxPayload, orderID)
	if err != nil {
		return nil, err
	}

	return p, nil
}
