package readers

import (
	"log"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/internal/pkg/sandbox"
	"github.com/FTChinese/superyard/pkg"
	"gorm.io/gorm"
)

func (env Env) countTestUser() (int64, error) {
	var count int64

	err := env.gormDBs.Read.
		Model(&sandbox.TestAccount{}).
		Count(&count).
		Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listTestUser(p gorest.Pagination) ([]sandbox.TestAccount, error) {
	var accounts = make([]sandbox.TestAccount, 0)

	err := env.gormDBs.Read.
		Limit(int(p.Limit)).
		Offset(int(p.Offset())).
		Find(&accounts).
		Error

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

	err := env.gormDBs.Write.Create(&account).Error

	if err != nil {
		return err
	}

	return nil
}

func (env Env) DeleteTestAccount(a sandbox.TestAccount) error {

	err := env.gormDBs.Delete.Delete(&a).Error

	if err != nil {
		return err
	}

	return nil
}

func (env Env) LoadSandboxAccount(ftcID string) (sandbox.TestAccount, error) {
	var a sandbox.TestAccount
	err := env.gormDBs.Read.Where("ftc_id", ftcID).
		First(&a).
		Error

	if err != nil {
		return sandbox.TestAccount{}, err
	}

	return a, nil
}

const stmtUpdatePassword = `
UPDATE cmstmp01.userinfo
SET password = MD5(?),
	updated_utc = UTC_TIMESTAMP()
WHERE user_id = ?
LIMIT 1
`

func (env Env) ChangePassword(a sandbox.TestAccount) error {
	return env.gormDBs.Write.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(a).Error
		if err != nil {
			return err
		}

		err = tx.Raw(stmtUpdatePassword, a.ClearPassword, a.FtcID).Error

		if err != nil {
			return err
		}

		return nil
	})
}
