package readers

import (
	"gitlab.com/ftchinese/superyard/models/reader"
)

func (env Env) CreateMember(m reader.Membership) error {
	m.Normalize()

	_, err := env.DB.NamedExec(stmtInsertMember, m)

	if err != nil {
		return err
	}

	return nil
}

// RetrieveMember load membership data.
func (env Env) RetrieveMember(id string) (reader.Membership, error) {
	var m reader.Membership

	err := env.DB.Get(&m, selectMemberByID, id)

	if err != nil {
		logger.WithField("trace", "Env.RetrieveMember").Error(err)
		return m, err
	}

	m.Normalize()

	return m, nil
}

// UpdateMember updates membership.
func (env Env) UpdateMember(m reader.Membership) error {

	m.Normalize()

	_, err := env.DB.NamedExec(stmtUpdateMember, m)

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
	var m reader.Membership
	if err := tx.Get(&m, selectMemberByID, id); err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	// Take a snapshot
	snapshot := reader.NewMemberSnapshot(m, reader.SnapshotReasonDelete)
	_, err = tx.NamedExec(insertMemberSnapshot, snapshot)
	if err != nil {
		log.Error(err)
		_ = tx.Rollback()
		return err
	}

	// Delete it.
	_, err = tx.Exec(stmtDeleteMember, id)
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
