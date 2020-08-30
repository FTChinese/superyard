package readers

import (
	"github.com/FTChinese/superyard/pkg/reader"
)

func (env Env) CreateSandboxUser(account reader.SandboxUser) error {
	tx, err := env.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(reader.StmtInsertSandbox, account)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtCreateSandboxUser, account)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (env Env) ListSandboxUsers() ([]reader.SandboxUser, error) {
	var accounts = make([]reader.SandboxUser, 0)
	if err := env.DB.Select(&accounts, reader.StmtListSandboxUsers); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (env Env) sandboxUserInfo(ftcId string) (reader.SandboxUser, error) {
	var a reader.SandboxUser
	err := env.DB.Get(&a, reader.StmtSandboxUser, ftcId)
	if err != nil {
		return reader.SandboxUser{}, err
	}

	return a, nil
}

type sandboxUserResult struct {
	value reader.SandboxUser
	err   error
}

func (env Env) asyncSandboxUser(ftcID string) <-chan sandboxUserResult {
	c := make(chan sandboxUserResult)

	go func() {
		defer close(c)
		s, err := env.sandboxUserInfo(ftcID)

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

	return reader.SandboxAccount{
		SandboxUser: sResult.value,
		Membership:  mResult.value,
	}, nil
}

func (env Env) SandboxUserExists(id string) (bool, error) {
	var found bool
	err := env.DB.Get(&found, reader.StmtSandboxExists, id)
	if err != nil {
		return false, err
	}

	return found, nil
}
