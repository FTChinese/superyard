package customer

import (
	"github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

func (env Env) RetrieveOrder(id string) (reader.Order, error) {
	var order reader.Order

	err := env.DB.Get(&order, stmtSelectOrder, id)
	if err != nil {
		logger.WithField("trace", "Env.RetrieveOrder").Error(err)
		return order, err
	}

	return order, nil
}

// ListOrders retrieves a user's orders that are paid successfully.
func (env Env) ListOrders(ids reader.AccountID, p gorest.Pagination) ([]reader.Order, error) {

	var orders = make([]reader.Order, 0)

	err := env.DB.Select(
		&orders,
		stmtReaderOrders,
		ids.FtcID,
		ids.UnionID,
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListOrders").Error(err)
		return nil, err
	}

	return orders, nil
}
