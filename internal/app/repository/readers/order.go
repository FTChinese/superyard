package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/ids"
	"github.com/FTChinese/superyard/pkg/subs"
	"github.com/labstack/gommon/log"
)

func (env Env) countOrders(uid ids.UserIDs) (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, subs.StmtCountOrder, uid.BuildFindInSet())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listOrders(uid ids.UserIDs, p gorest.Pagination) ([]subs.Order, error) {
	var orders = make([]subs.Order, 0)

	err := env.dbs.Read.Select(
		&orders,
		subs.StmtListOrders,
		uid.BuildFindInSet(),
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
func (env Env) ListOrders(uid ids.UserIDs, p gorest.Pagination) (subs.OrderList, error) {
	countCh := make(chan int64)
	listCh := make(chan subs.OrderList)

	go func() {
		defer close(countCh)
		n, err := env.countOrders(uid)
		if err != nil {
			log.Error(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listOrders(uid, p)
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
