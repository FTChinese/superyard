package readers

import (
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) CreateSandboxUser(account reader.SandboxFtcAccount) error {
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
func (env Env) sandboxJoinedSchema(ftcId string) (reader.SandboxJoinedAccountSchema, error) {
	var a reader.SandboxJoinedAccountSchema
	err := env.db.Get(&a, reader.StmtSandboxJoinedAccount, ftcId)
	if err != nil {
		return reader.SandboxJoinedAccountSchema{}, err
	}

	return a, nil
}

type sandboxUserResult struct {
	value reader.SandboxJoinedAccountSchema
	err   error
}

func (env Env) asyncSandboxUser(ftcID string) <-chan sandboxUserResult {
	c := make(chan sandboxUserResult)

	go func() {
		defer close(c)
		s, err := env.sandboxJoinedSchema(ftcID)

		c <- sandboxUserResult{
			value: s,
			err:   err,
		}
	}()

	return c
}

func (env Env) LoadSandboxAccount(ftcID string) (reader.SandboxAccount, error) {
	sChan, mChan := env.asyncSandboxUser(ftcID), env.asyncMembership(ftcID)

	sResult, mResult := <-sChan, <-mChan

	if sResult.err != nil {
		return reader.SandboxAccount{}, sResult.err
	}

	if mResult.err != nil {
		return reader.SandboxAccount{}, mResult.err
	}

	return sResult.value.Build(mResult.value), nil
}

func (env Env) SandboxUserExists(id string) (bool, error) {
	var found bool
	err := env.db.Get(&found, reader.StmtSandboxExists, id)
	if err != nil {
		return false, err
	}

	return found, nil
}

func (env Env) ChangePassword(u reader.SandboxFtcAccount) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(reader.StmtUpdateClearPassword, u)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtUpdatePassword, u)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
