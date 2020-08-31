package readers

import (
	"database/sql"
	"github.com/FTChinese/go-rest/enum"
	"github.com/FTChinese/superyard/pkg/paywall"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
	"go.uber.org/zap"
)

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

// UpdateMember changes a membership directly.
func (env Env) UpdateMember(input reader.MemberInput, plan paywall.Plan) (subs.ConfirmationResult, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	tx, err := env.BeginMemberTx()
	if err != nil {
		sugar.Error(err)
		return subs.ConfirmationResult{}, err
	}

	// Retrieve current membership
	current, err := tx.RetrieveMember(
		input.CompoundID)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}
	current = current.Normalize()

	m := current.Update(input, plan)

	sugar.Infof("Updated membership %+v", m)

	// Update it.
	err = tx.UpdateMember(m)
	if err != nil {
		sugar.Error(err)
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

func (env Env) DeleteMember(compoundID string) (reader.MemberSnapshot, error) {
	tx, err := env.BeginMemberTx()
	if err != nil {
		return reader.MemberSnapshot{}, err
	}

	m, err := tx.RetrieveMember(compoundID)
	if err != nil {
		_ = tx.Rollback()
		return reader.MemberSnapshot{}, err
	}
	m = m.Normalize()

	err = tx.DeleteMember(m.CompoundID.String)
	if err != nil {
		_ = tx.Rollback()
		return reader.MemberSnapshot{}, err
	}

	if err := tx.Commit(); err != nil {
		return reader.MemberSnapshot{}, err
	}

	return reader.NewSnapshot(enum.SnapshotReasonDelete, m), nil
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
