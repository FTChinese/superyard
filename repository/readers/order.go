package readers

import (
	"database/sql"
	"github.com/FTChinese/go-rest"
	"gitlab.com/ftchinese/superyard/pkg/subs"
)

// ListOrders retrieves a user's orders.
func (env Env) ListOrders(ids subs.CompoundIDs, p gorest.Pagination) ([]subs.Order, error) {

	var orders = make([]subs.Order, 0)

	err := env.DB.Select(
		&orders,
		subs.StmtListOrders,
		ids.BuildFindInSet(),
		p.Limit,
		p.Offset())

	if err != nil {
		logger.WithField("trace", "Env.ListOrders").Error(err)
		return nil, err
	}

	return orders, nil
}

// RetrieveOrder retrieves a single order by trade_no column.
func (env Env) RetrieveOrder(id string) (subs.Order, error) {
	var order subs.Order

	err := env.DB.Get(&order, subs.StmtSelectOrder, id)
	if err != nil {
		logger.WithField("trace", "Env.RetrieveOrder").Error(err)
		return order, err
	}

	return order, nil
}

// ConfirmOrder is used to confirmed an order.
// Errors returned:
// subs.ErrAlreadyConfirmed
// subs.ErrAlreadyUpgraded
func (env Env) ConfirmOrder(id string) (subs.ConfirmationResult, error) {
	log := logger.WithField("trace", "Env.ConfirmOrder")

	tx, err := env.DB.Beginx()
	if err != nil {
		log.Error(err)
		return subs.ConfirmationResult{}, err
	}

	var order subs.Order
	if err := tx.Get(&order, subs.StmtSelectOrder, id); err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}
	log.Infof("Order retrieved: %s", order.ID)

	// Retrieve membership. sql.ErrNoRows should be treated
	// as valid.
	var member subs.Membership
	err = tx.Get(&member, subs.StmtMembership, order.CompoundID)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}
	member = member.Normalize()

	builder := subs.NewConfirmationBuilder(order, member)

	if err := builder.Validate(); err != nil {
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	result, err := builder.Build()
	if err != nil {
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	// Save the confirmed order
	_, err = tx.NamedExec(subs.StmtConfirmOrder, result.Order)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	var stmtUpsertMember string
	if member.IsZero() {
		stmtUpsertMember = subs.StmtInsertMember
	} else {
		stmtUpsertMember = subs.StmtUpdateMember
	}
	_, err = tx.NamedExec(stmtUpsertMember, result.Membership)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	// If old membership is not empty, back up it.
	if !result.Snapshot.IsZero() {
		_, err = tx.NamedExec(subs.InsertMemberSnapshot, result.Snapshot)
		if err != nil {
			log.Error(err)
			_ = tx.Rollback()
			return subs.ConfirmationResult{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		return subs.ConfirmationResult{}, err
	}

	log.Infof("Confirmed order finished")

	return result, nil
}
