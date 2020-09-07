package readers

import (
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) CreateSandboxUser(account reader.FtcAccount) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(reader.StmtInsertSandbox, account)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtCreateReader, account)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtCreateProfile, account)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (env Env) DeleteSandboxAccount(id string) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(reader.StmtDeleteSandbox, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtDeleteAccount, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtDeleteProfile, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtDeleteMember, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (env Env) ListSandboxFtcAccount() ([]reader.FtcAccount, error) {
	var accounts = make([]reader.FtcAccount, 0)
	if err := env.db.Select(&accounts, reader.StmtListSandboxUsers); err != nil {
		return nil, err
	}

	return accounts, nil
}

// retrieves sandbox user's ftc account + wechat
func (env Env) sandboxJoinedSchema(ftcId string) (reader.JoinedAccountSchema, error) {
	var a reader.JoinedAccountSchema
	err := env.db.Get(&a, reader.StmtSandboxJoinedAccount, ftcId)
	if err != nil {
		return reader.JoinedAccountSchema{}, err
	}

	return a, nil
}

func (env Env) asyncSandboxJoinedAccount(ftcID string) <-chan accountAsyncResult {
	c := make(chan accountAsyncResult)

	go func() {
		defer close(c)
		s, err := env.sandboxJoinedSchema(ftcID)

		c <- accountAsyncResult{
			value: s,
			err:   err,
		}
	}()

	return c
}

func (env Env) LoadSandboxAccount(ftcID string) (reader.Account, error) {
	aChan, mChan := env.asyncSandboxJoinedAccount(ftcID), env.asyncMembership(ftcID)

	aResult, mResult := <-aChan, <-mChan

	if aResult.err != nil {
		return reader.Account{}, aResult.err
	}

	if mResult.err != nil {
		return reader.Account{}, mResult.err
	}

	return aResult.value.BuildAccount(mResult.value), nil
}

func (env Env) SandboxUserExists(id string) (bool, error) {
	var found bool
	err := env.db.Get(&found, reader.StmtSandboxExists, id)
	if err != nil {
		return false, err
	}

	return found, nil
}

func (env Env) ChangePassword(s reader.SandboxPasswordUpdater) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(reader.StmtUpdateClearPassword, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtUpdatePassword, s)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
