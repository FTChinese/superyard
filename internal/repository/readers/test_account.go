package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/reader"
	"log"
)

func (env Env) countTestUser() (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, reader.StmtCountTestUser)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listTestUser(p gorest.Pagination) ([]reader.FtcAccount, error) {
	var accounts = make([]reader.FtcAccount, 0)
	err := env.dbs.Read.Select(
		&accounts,
		reader.StmtListTestUsers,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (env Env) ListTestFtcAccount(p gorest.Pagination) (reader.FtcAccountList, error) {
	countCh := make(chan int64)
	listCh := make(chan reader.FtcAccountList)

	go func() {
		defer close(countCh)
		n, err := env.countTestUser()
		if err != nil {
			log.Print(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listTestUser(p)
		listCh <- reader.FtcAccountList{
			Total:      0,
			Pagination: gorest.Pagination{},
			Data:       list,
			Err:        err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return reader.FtcAccountList{}, listResult.Err
	}

	return reader.FtcAccountList{
		Total:      count,
		Pagination: p,
		Data:       listResult.Data,
		Err:        nil,
	}, nil
}

func (env Env) CreateTestUser(account reader.FtcAccount) error {
	tx, err := env.dbs.Write.Beginx()
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
	tx, err := env.dbs.Delete.Beginx()
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

// retrieves sandbox user's ftc account + wechat
func (env Env) testJoinedSchema(ftcId string) (reader.JoinedAccountSchema, error) {
	var a reader.JoinedAccountSchema
	err := env.dbs.Read.Get(&a, reader.StmtTestJoinedAccount, ftcId)
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
	err := env.dbs.Read.Get(&found, reader.StmtTestUserExists, id)
	if err != nil {
		return false, err
	}

	return found, nil
}

func (env Env) ChangePassword(s reader.TestPasswordUpdater) error {
	tx, err := env.dbs.Write.Beginx()
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
