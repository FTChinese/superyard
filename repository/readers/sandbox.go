package readers

import "github.com/FTChinese/superyard/pkg/reader"

func (env Env) CreateSandbox(account reader.SandboxAccount) error {
	tx, err := env.DB.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(reader.StmtInsertSandbox, account)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(reader.StmtCreateAccount, account)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (env Env) RetrieveSandboxReader(id string) (reader.SandboxAccount, error) {
	var a reader.SandboxAccount
	err := env.DB.Get(&a, reader.StmtSandboxAccount, id)
	if err != nil {
		return reader.SandboxAccount{}, err
	}

	return a, nil
}

func (env Env) SandboxAccountExists(id string) (bool, error) {
	var found bool
	err := env.DB.Get(&found, reader.StmtSandboxExists, id)
	if err != nil {
		return false, err
	}

	return found, nil
}
