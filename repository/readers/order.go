package readers

import (
	"database/sql"
	"errors"
	"github.com/FTChinese/go-rest"
	"github.com/FTChinese/go-rest/enum"
	"gitlab.com/ftchinese/superyard/models/builder"
	"gitlab.com/ftchinese/superyard/models/reader"
	"time"
)

// ListOrders retrieves a user's orders.
func (env Env) ListOrders(ids reader.AccountID, p gorest.Pagination) ([]reader.Order, error) {

	inBuilder := builder.
		NewInBuilder(ids.QueryArgs()...).
		Append(p.Limit, p.Offset())

	var orders = make([]reader.Order, 0)

	err := env.DB.Select(
		&orders,
		stmtListOrders(inBuilder.PlaceHolder()),
		inBuilder.Values()...)

	if err != nil {
		logger.WithField("trace", "Env.ListOrders").Error(err)
		return nil, err
	}

	return orders, nil
}

// CreateOrder inserts an new order record.
func (env Env) CreateOrder(order reader.Order) error {
	_, err := env.DB.NamedExec(stmtInsertOrder, order)
	if err != nil {
		return err
	}

	return nil
}

// RetrieveOrder retrieves a single order by trade_no column.
func (env Env) RetrieveOrder(id string) (reader.Order, error) {
	var order reader.Order

	err := env.DB.Get(&order, stmtAnOrder, id)
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

	var order reader.Order
	if err := tx.Get(&order, stmtAnOrder, id); err != nil {
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
	var member reader.Membership
	if err := tx.Get(&member, memberByCompoundID, order.CompoundID); err != nil && err != sql.ErrNoRows {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}
	if member.ID.IsZero() {
		member.GenerateID()
	}

	log.Infof("Member retrieved: %+v", member)

	// Cannot upgrade a premium meber
	if order.Kind == reader.SubsKindUpgrade && member.Tier == enum.TierPremium {
		log.Infof("Order %s is trying to upgrade a premium member %s", order.ID, member.ID.String)
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
	_, err = tx.NamedExec(stmtConfirmOrder, confirmedOrder)
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
	newMember.Normalize()

	log.Infof("New membership created")

	if member.IsZero() {
		_, err := tx.NamedExec(stmtInsertMember, newMember)
		if err != nil {
			log.Error(err)
			_ = tx.Rollback()
			return err
		}
	} else {
		_, err := tx.NamedExec(stmtUpdateMember, newMember)
		if err != nil {
			log.Error(err)
			_ = tx.Rollback()
			return err
		}
	}

	// If old membership is not empty, back up it.
	if !member.IsZero() {
		snapshot := reader.NewMemberSnapshot(member, order.Kind.SnapshotReason())
		_, err = tx.NamedExec(insertMemberSnapshot, snapshot)
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
