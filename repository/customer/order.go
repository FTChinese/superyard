package customer

import (
	"github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/backyard-api/models/reader"
)

// ListOrders retrieves a user's orders that are paid successfully.
func (env Env) ListOrders(ids reader.AccountID, p gorest.Pagination) ([]reader.Order, error) {

	var orders = make([]reader.Order, 0)

	err := env.DB.Select(
		&orders,
		stmtListOrders,
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

func (env Env) RetrieveOrder(id string) (reader.Order, error) {
	var order reader.Order

	err := env.DB.Get(&order, stmtSelectOneOrder, id)
	if err != nil {
		logger.WithField("trace", "Env.RetrieveOrder").Error(err)
		return order, err
	}

	return order, nil
}

// CreateOrder inserts an new order record.
func (env Env) CreateOrder(order reader.Order) error {
	_, err := env.DB.NamedExec(stmtCreateOrder, order)

	if err != nil {
		logger.WithField("trace", "Env.CreateOrder").Error(err)

		return err
	}

	return nil
}

// UpdateOrder is used to confirmed an order.
func (env Env) UpdateOrder(order reader.Order) error {
	_, err := env.DB.NamedExec(stmtCreateOrder, order)

	if err != nil {
		logger.WithField("trace", "Env.UpdateOrder").Error(err)
		return err
	}

	return nil
}
