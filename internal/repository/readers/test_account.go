package readers

import (
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) CreateTestUser(account reader.FtcAccount) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(reader.StmtInsertTestAccount, account)
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

func (env Env) DeleteTestAccount(id string) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.Exec(reader.StmtDeleteTestUser, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(reader.StmtDeleteAccount, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(reader.StmtDeleteProfile, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(reader.StmtDeleteMember, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (env Env) ListTestFtcAccount() ([]reader.FtcAccount, error) {
	var accounts = make([]reader.FtcAccount, 0)
	if err := env.db.Select(&accounts, reader.StmtListTestUsers); err != nil {
		return nil, err
	}

	return accounts, nil
}

// retrieves sandbox user's ftc account + wechat
func (env Env) testJoinedSchema(ftcId string) (reader.JoinedAccountSchema, error) {
	var a reader.JoinedAccountSchema
	err := env.db.Get(&a, reader.StmtTestJoinedAccount, ftcId)
	if err != nil {
		return reader.JoinedAccountSchema{}, err
	}

	return a, nil
}

func (env Env) asyncSandboxJoinedAccount(ftcID string) <-chan accountAsyncResult {
	c := make(chan accountAsyncResult)

	go func() {
		defer close(c)
		s, err := env.testJoinedSchema(ftcID)

		c <- accountAsyncResult{
			value: s,
			err:   err,
		}
	}()

	return c
}

func (env Env) LoadSandboxAccount(ftcID string) (reader.Account, error) {
	aChan, mChan := env.asyncSandboxJoinedAccount(ftcID), env.asyncAccountMember(ftcID)

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
	err := env.db.Get(&found, reader.StmtTestUserExists, id)
	if err != nil {
		return false, err
	}

	return found, nil
}

func (env Env) ChangePassword(s reader.TestPasswordUpdater) error {
	tx, err := env.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(reader.StmtUpdateTestUserPassword, s)
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
