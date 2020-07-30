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
// 3 groups of data are involved:
// * The new Membership;
// * Current membership from db;
// * Snapshot based on current membership.
func (env Env) UpdateMember(m subs.Membership, creator string) error {
	m = m.Normalize()

	tx, err := env.DB.Beginx()
	if err != nil {
		logger.WithField("trace", "UpdateMember.Beginx").Error(err)
		return err
	}

	// Retrieve the membership
	var current subs.Membership
	if err := tx.Get(&current, subs.StmtMembership, m.CompoundID); err != nil {
		logger.WithField("trace", "UpdateMember.RetrieveMember").Error(err)
		_ = tx.Rollback()
		return err
	}
	current.Normalize()

	// Take a snapshot
	snapshot := current.Snapshot(enum.SnapshotReasonManual).
		WithCreator(creator)
	_, err = tx.NamedExec(subs.InsertMemberSnapshot, snapshot)
	if err != nil {
		logger.WithField("trace", "UpdateMember.Snapshot").Error(err)
		_ = tx.Rollback()
		return err
	}

	// Update it.
	_, err = tx.NamedExec(subs.StmtUpdateMember, m)
	if err != nil {
		logger.WithField("trace", "UpdateMember.Update").Error(err)
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.WithField("trace", "UpdateMember.Commit").Error(err)
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
