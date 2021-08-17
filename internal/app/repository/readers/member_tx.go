package readers

import (
	"database/sql"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/jmoiron/sqlx"
)

type MemberTx struct {
	*sqlx.Tx
}

func NewMemberTx(tx *sqlx.Tx) MemberTx {
	return MemberTx{
		Tx: tx,
	}
}

// RetrieveMember retrieves a user's membership by a compound id, which might be ftc id or union id.
// Use SQL FIND_IN_SET(compoundId, vip_id, vip) to verify it against two columns.
// Returns zero value of membership if not found.
func (tx MemberTx) RetrieveMember(compoundID string) (reader.Membership, error) {
	var m reader.Membership

	err := tx.Get(
		&m,
		reader.StmtLockMember,
		compoundID,
	)

	if err != nil && err != sql.ErrNoRows {
		return m, err
	}

	// Treat a non-existing member as a valid value.
	return m.Normalize(), nil
}

// CreateMember creates a new membership for an order of a new subscription.
func (tx MemberTx) CreateMember(m reader.Membership) error {
	m = m.Normalize()

	_, err := tx.NamedExec(
		reader.StmtCreateMember,
		m,
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateMember updates existing membership for orders whose kind is renew or upgrade.
func (tx MemberTx) UpdateMember(m reader.Membership) error {
	m = m.Normalize()

	_, err := tx.NamedExec(
		reader.StmtUpdateMember,
		m)

	if err != nil {
		return err
	}

	return nil
}

func (tx MemberTx) DeleteMember(compoundID string) error {
	_, err := tx.Exec(
		reader.StmtDeleteMember,
		compoundID)

	if err != nil {
		return err
	}

	return nil
}
