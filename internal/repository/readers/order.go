package readers

import (
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/subs"
	"go.uber.org/zap"
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
func (env Env) RetrieveOrder(id string) (subs.Order, error) {
	var order subs.Order

	err := env.db.Get(&order, subs.StmtOrder, id)
	if err != nil {
		return order, err
	}

	return order, nil
}

// ConfirmOrder is used to confirmed an order.
// Errors returned:
// subs.ErrAlreadyConfirmed
// subs.ErrAlreadyUpgraded
func (env Env) ConfirmOrder(id string) (subs.ConfirmationResult, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	tx, err := env.BeginMemberTx()
	if err != nil {
		sugar.Error(err)
		return subs.ConfirmationResult{}, err
	}

	order, err := tx.RetrieveOrder(id)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}
	sugar.Infof("Order retrieved: %s", order.ID)

	// Retrieve membership. sql.ErrNoRows should be treated
	// as valid.
	member, err := tx.RetrieveMember(order.CompoundID)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}
	member = member.Normalize()

	builder := subs.NewConfirmationBuilder(order, member)

	if err := builder.Validate(); err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	result, err := builder.Build()
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	// Saved confirmed order.
	err = tx.ConfirmOrder(result.Order)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	// Flag upgrade balance source as consumed.
	if result.Order.Kind == enum.OrderKindUpgrade {
		err := tx.ProratedOrdersUsed(result.Order.ID)

		if err != nil {
			sugar.Error(err)
			_ = tx.Rollback()
			return subs.ConfirmationResult{}, err
		}
	}

	if member.IsZero() {
		if err := tx.CreateMember(result.Membership); err != nil {
			sugar.Error(err)
			_ = tx.Rollback()

			return subs.ConfirmationResult{}, err
		}
	} else {
		if err := tx.UpdateMember(result.Membership); err != nil {
			sugar.Error(err)
			_ = tx.Rollback()

			return subs.ConfirmationResult{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		sugar.Error(err)
		return subs.ConfirmationResult{}, err
	}

	sugar.Infof("Confirmed order finished")

	return result, nil
}
