package readers

import (
	"database/sql"
	"github.com/FTChinese/superyard/pkg/reader"
	"github.com/FTChinese/superyard/pkg/subs"
)

// RetrieveMember load membership data.
// The id might a ftc uuid or wechat union id.
func (env Env) MemberByCompoundID(compoundID string) (reader.Membership, error) {
	var m reader.Membership

	err := env.dbs.Read.Get(&m, reader.StmtSelectMember, compoundID)

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

func (env Env) asyncAccountMember(compoundID string) <-chan memberAsyncResult {
	c := make(chan memberAsyncResult)

	go func() {
		m, err := env.MemberByCompoundID(compoundID)

		c <- memberAsyncResult{
			value: m,
			err:   err,
		}
	}()

	return c
}

// CreateFtcMember creates membership purchased via ali or wx pay for an account.
// If the account is not found, or membership already exists,
// error will be returned.
func (env Env) CreateFtcMember(input subs.FtcSubsCreationInput) (reader.Account, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	// Find user's account first. Stop if not found.
	a, err := env.JoinedAccountByFtcOrWx(input.IDs)
	if err != nil {
		return reader.Account{}, err
	}

	tx, err := env.BeginMemberTx()
	if err != nil {
		sugar.Error(err)
		return reader.Account{}, err
	}

	// Then check if this account has membership. We should stop if membership present.
	current, err := tx.RetrieveMember(
		a.MustGetCompoundID())
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return reader.Account{}, err
	}
	if ve := current.ValidateCreateFtc(); ve != nil {
		return reader.Account{}, ve
	}

	// If account not found, then membership should not be present.
	newMmb := subs.ManualCreateMember(a, input)

	err = tx.CreateMember(newMmb)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return reader.Account{}, err
	}

	if err := tx.Commit(); err != nil {
		return reader.Account{}, err
	}

	return reader.Account{
		JoinedAccount: a,
		Membership:    newMmb,
	}, nil
}

// UpdateFtcMember changes an ftc membership directly.
func (env Env) UpdateFtcMember(compoundID string, input subs.FtcSubsUpdateInput) (subs.ConfirmationResult, error) {
	defer env.logger.Sync()
	sugar := env.logger.Sugar()

	tx, err := env.BeginMemberTx()
	if err != nil {
		sugar.Error(err)
		return subs.ConfirmationResult{}, err
	}

	// Lock and retrieve membership. If not found, we got noting to update.
	current, err := tx.RetrieveMember(compoundID)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}
	current = current.Normalize()

	// Check whether current membership permits updating.
	if err := current.ValidateUpdateFtc(); err != nil {
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	newMmb := subs.ManualUpdateMember(current, input)

	sugar.Infof("Updated membership %+v", newMmb)

	err = tx.UpdateMember(newMmb)
	if err != nil {
		sugar.Error(err)
		_ = tx.Rollback()
		return subs.ConfirmationResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return subs.ConfirmationResult{}, err
	}

	return subs.ConfirmationResult{
		Membership: newMmb,
		Snapshot:   reader.NewSnapshot(current),
	}, nil
}

func (env Env) DeleteFtcMember(compoundID string) (reader.MemberSnapshot, error) {
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

	// Only membership purchased via ali or wx is allowed to be deleted directly.
	if m.IsZero() || !m.IsAliOrWxPay() {
		_ = tx.Rollback()
		return reader.MemberSnapshot{}, sql.ErrNoRows
	}

	err = tx.DeleteMember(m.CompoundID.String)
	if err != nil {
		_ = tx.Rollback()
		return reader.MemberSnapshot{}, err
	}

	if err := tx.Commit(); err != nil {
		return reader.MemberSnapshot{}, err
	}

	return reader.NewSnapshot(m), nil
}
