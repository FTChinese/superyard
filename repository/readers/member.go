package readers

import (
	"github.com/FTChinese/go-rest/enum"
	"gitlab.com/ftchinese/superyard/pkg/subs"
)

func (env Env) CreateMember(m subs.Membership) error {
	m = m.Normalize()

	_, err := env.DB.NamedExec(subs.StmtInsertMember, m)

	if err != nil {
		return err
	}

	return nil
}

// RetrieveMember load membership data.
func (env Env) RetrieveMember(id string) (subs.Membership, error) {
	var m subs.Membership

	err := env.DB.Get(&m, subs.StmtMembership, id)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveMember").Error(err)
		return m, err
	}

	return m.Normalize(), nil
}

// UpdateMember updates membership.
func (env Env) UpdateMember(m subs.Membership) error {

	m = m.Normalize()

	_, err := env.DB.NamedExec(subs.StmtUpdateMember, m)

	if err != nil {
		logger.WithField("trace", "Env.UpdateMember").Error(err)
		return err
	}

	return nil
}

// DeleteMember removes the specified record from ftc_vip
// and backup it.
func (env Env) DeleteMember(id string) error {
	log := logger.WithField("trace", "Env.DeleteMember")

	tx, err := env.DB.Beginx()
	if err != nil {
		log.Error(err)
		return err
	}

	// Retrieve the membership
	var m subs.Membership
	if err := tx.Get(&m, subs.StmtMembership, id); err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	// Take a snapshot
	snapshot := m.Snapshot(enum.SnapshotReasonDelete)
	_, err = tx.NamedExec(subs.InsertMemberSnapshot, snapshot)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	// Delete it.
	_, err = tx.Exec(subs.StmtDeleteMember, id)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// FindMemberForOrder tries to find the current membership
// by an order's compound id, which might be either
// ftc id or wechat union id.
func (env Env) FindMemberForOrder(ftcOrUnionID string) (subs.Membership, error) {
	var m subs.Membership

	err := env.DB.Get(&m, subs.StmtMemberForOrder, ftcOrUnionID)

	if err != nil {
		logger.WithField("trace", "Env.FindMemberForOrder").Error(err)

		return m, err
	}

	return m.Normalize(), nil
}
