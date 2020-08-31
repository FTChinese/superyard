package readers

import (
	"database/sql"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
)

func (env Env) CreateMember(m reader.Membership) error {
	m = m.Normalize()

	_, err := env.db.NamedExec(reader.StmtInsertMember, m)

	if err != nil {
		return err
	}

	return nil
}

// RetrieveMember load membership data.
// The id might a ftc uuid or wechat union id.
func (env Env) RetrieveMember(id string) (reader.Membership, error) {
	var m reader.Membership

	err := env.db.Get(&m, reader.StmtMembership, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return reader.Membership{}, nil
		}
		return reader.Membership{}, err
	}

	return m.Normalize(), nil
}

type memberAsyncResult struct {
	value reader.Membership
	err   error
}

func (env Env) asyncMembership(id string) <-chan memberAsyncResult {
	c := make(chan memberAsyncResult)

	go func() {
		m, err := env.RetrieveMember(id)

		c <- memberAsyncResult{
			value: m,
			err:   err,
		}
	}()

	return c
}

// UpdateMember updates membership.
// 3 groups of data are involved:
// * The new Membership;
// * Current membership from db;
// * Snapshot based on current membership.
func (env Env) UpdateMember(input reader.MemberInput, plan paywall.Plan) (subs.ConfirmationResult, error) {

	tx, err := env.db.Beginx()
	if err != nil {
		logger.WithField("trace", "UpdateMember.Beginx").Error(err)
		return subs.ConfirmationResult{}, err
	}

	// Retrieve current membership
	var current reader.Membership
	err = tx.Get(
		&current,
		reader.StmtMembership,
		input.CompoundID)
	if err != nil {
		logger.WithField("trace", "UpdateMember.RetrieveMember").Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}
	current.Normalize()

	m := current.Update(input, plan).Normalize()

	// Update it.
	_, err = tx.NamedExec(
		reader.StmtUpdateMember,
		m)
	if err != nil {
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return subs.ConfirmationResult{}, err
	}

	return subs.ConfirmationResult{
		Membership: m,
		Snapshot: reader.NewSnapshot(
			enum.SnapshotReasonManual,
			current),
	}, nil
}

func (env Env) SnapshotMember(s reader.MemberSnapshot) error {
	_, err := env.db.NamedExec(
		reader.InsertMemberSnapshot,
		s)
	if err != nil {
		return err
	}

	return nil
}
