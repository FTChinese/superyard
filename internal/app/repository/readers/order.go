package readers

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/labstack/gommon/log"
)

func (env Env) countOrders(ids reader.IDs) (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, subs.StmtCountOrder, ids.BuildFindInSet())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listOrders(ids reader.IDs, p gorest.Pagination) ([]subs.Order, error) {
	var orders = make([]subs.Order, 0)

	err := env.dbs.Read.Select(
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

// ListOrders retrieves a user's orders.
// Turn reader's possible into a format used in
// MySQL function FIND_IN_SET.
func (env Env) ListOrders(ids reader.IDs, p gorest.Pagination) (subs.OrderList, error) {
	countCh := make(chan int64)
	listCh := make(chan subs.OrderList)

	go func() {
		defer close(countCh)
		n, err := env.countOrders(ids)
		if err != nil {
			log.Error(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listOrders(ids, p)
		listCh <- subs.OrderList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return subs.OrderList{}, listResult.Err
	}

	listResult.Total = count
	listResult.Pagination = p

	return listResult, nil
}

// RetrieveOrder retrieves a single order by trade_no column.
func (env Env) RetrieveOrder(orderID string) (subs.Order, error) {
	var order subs.Order

	err := env.dbs.Read.Get(&order, subs.StmtOrder, orderID)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (env Env) AliWebhook(orderID string) ([]subs.AliPayload, error) {
	var p = make([]subs.AliPayload, 0)

	err := env.dbs.Read.Select(&p, subs.StmtAliPayload, orderID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (env Env) WxWebhook(orderID string) ([]subs.WxPayload, error) {
	var p = make([]subs.WxPayload, 0)

	err := env.dbs.Read.Select(&p, subs.StmtWxPayload, orderID)
	if err != nil {
		return nil, err
	}

	return p, nil
}
