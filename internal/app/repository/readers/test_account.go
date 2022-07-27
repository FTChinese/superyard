package readers

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/internal/pkg/sandbox"
	"github.com/FTChinese/superyard/pkg"
	"log"
)

func (env Env) countTestUser() (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, sandbox.StmtCountTestUser)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listTestUser(p gorest.Pagination) ([]sandbox.TestAccount, error) {
	var accounts = make([]sandbox.TestAccount, 0)
	err := env.dbs.Read.Select(
		&accounts,
		sandbox.StmtListTestUsers,
		p.Limit,
		p.Offset())

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (env Env) ListTestFtcAccount(p gorest.Pagination) (pkg.PagedList[sandbox.TestAccount], error) {
	countCh := make(chan int64)
	listCh := make(chan pkg.AsyncResult[[]sandbox.TestAccount])

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
		listCh <- pkg.AsyncResult[[]sandbox.TestAccount]{
			Value: list,
			Err:   err,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return pkg.PagedList[sandbox.TestAccount]{}, listResult.Err
	}

	return pkg.PagedList[sandbox.TestAccount]{
		Total:      count,
		Pagination: p,
		Data:       listResult.Value,
	}, nil
}

func (env Env) CreateTestUser(account sandbox.TestAccount) error {

	_, err := env.dbs.Write.NamedExec(
		sandbox.StmtInsertTestAccount,
		account)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) DeleteTestAccount(id string) error {

	_, err := env.dbs.Write.Exec(sandbox.StmtDeleteTestUser, id)
	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadSandboxAccount(ftcID string) (sandbox.TestAccount, error) {
	var a sandbox.TestAccount
	err := env.dbs.Read.Get(
		&a,
		sandbox.StmtRetrieveTestUser,
		ftcID)

	if err != nil {
		return sandbox.TestAccount{}, err
	}

	return a, nil
}

func (env Env) ChangePassword(a sandbox.TestAccount) error {
	tx, err := env.dbs.Write.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec(sandbox.StmtUpdateTestUserPassword, a)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(sandbox.StmtUpdatePassword, a)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
