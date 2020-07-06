package readers

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/enum"
	"gitlab.com/ftchinese/superyard/pkg/subs"
	"time"
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
func (env Env) ConfirmOrder(id string) error {
	log := logger.WithField("trace", "Env.ConfirmOrder")

	tx, err := env.DB.Beginx()
	if err != nil {
		log.Error(err)
		return err
	}

	var order subs.Order
	if err := tx.Get(&order, subs.StmtSelectOrder, id); err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}
	log.Infof("Order retrieved: %s", order.ID)

	// If order is already confirmed.
	if order.IsConfirmed() {
		_ = tx.Rollback()
		return errors.New("order already confirmed")
	}

	// Retrieve membership. sql.ErrNoRows should be treated
	// as valid.
	var member subs.Membership
	err = tx.Get(&member, subs.StmtMembership, order.CompoundID)
	if err != nil && err != sql.ErrNoRows {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	log.Infof("Member retrieved: %+v", member)

	// Cannot upgrade a premium member
	if order.Kind == subs.KindUpgrade && member.Tier == enum.TierPremium {
		log.Infof("Order %s is trying to upgrade a premium member", order.ID)
		_ = tx.Rollback()
		return errors.New("cannot upgrade a premium membership")
	}

	// Create the confirmed order
	confirmedOrder, err := order.Confirm(member, time.Now())
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}
	log.Info("Order confirmed")

	// Save the confirmed order
	_, err = tx.NamedExec(subs.StmtConfirmOrder, confirmedOrder)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	// Create new membership based on the confirmed order
	newMember, err := member.FromAliOrWx(confirmedOrder)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}
	// This step is important to keep compatibility.
	newMember = newMember.Normalize()

	log.Infof("New membership created")

	if member.IsZero() {
		_, err := tx.NamedExec(subs.StmtInsertMember, newMember)
		if err != nil {
			log.Error(err)
			_ = tx.Rollback()
			return err
		}
	} else {
		_, err := tx.NamedExec(subs.StmtUpdateMember, newMember)
		if err != nil {
			log.Error(err)
			_ = tx.Rollback()
			return err
		}
	}

	// If old membership is not empty, back up it.
	if !member.IsZero() {
		snapshot := member.Snapshot(order.Kind.SnapshotReason())
		_, err = tx.NamedExec(subs.InsertMemberSnapshot, snapshot)
		if err != nil {
			log.Error(err)
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		return err
	}

	log.Infof("Confirmed order finished")

	return nil
}
