package readers

import (
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) CreateMember(m reader.Membership) error {
	m = m.Normalize()

	_, err := env.DB.NamedExec(reader.StmtInsertMember, m)

	if err != nil {
		return err
	}

	return nil
}

// RetrieveMember load membership data.
func (env Env) RetrieveMember(id string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, reader.StmtMembership, id)

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
func (env Env) UpdateMember(m reader.Membership, creator string) error {
	m = m.Normalize()

	tx, err := env.DB.Beginx()
	if err != nil {
		logger.WithField("trace", "UpdateMember.Beginx").Error(err)
		return err
	}

	// Retrieve the membership
	var current reader.Membership
	if err := tx.Get(&current, reader.StmtMembership, m.CompoundID); err != nil {
		logger.WithField("trace", "UpdateMember.RetrieveMember").Error(err)
		_ = tx.Rollback()
		return err
	}
	current.Normalize()

	// Take a snapshot
	snapshot := reader.NewSnapshot(enum.SnapshotReasonManual, current).WithCreator(creator)
	if !snapshot.IsZero() {
		_, err = tx.NamedExec(reader.InsertMemberSnapshot, snapshot)
		if err != nil {
			logger.WithField("trace", "UpdateMember.Snapshot").Error(err)
			_ = tx.Rollback()
			return err
		}
	}

	// Update it.
	_, err = tx.NamedExec(reader.StmtUpdateMember, m)
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
func (env Env) FindMemberForOrder(ftcOrUnionID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, reader.StmtMemberForOrder, ftcOrUnionID)

	if err != nil {
		logger.WithField("trace", "Env.FindMemberForOrder").Error(err)

		return m, err
	}

	return m.Normalize(), nil
}
