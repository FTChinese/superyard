package admin

import (
	"log"

	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/internal/pkg/user"
	"github.com/FTChinese/superyard/pkg"
)

func (env Env) countStaff() (int64, error) {
	var count int64

	err := env.dbs.Read.Get(&count, user.StmtCountStaff)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listStaff(p gorest.Pagination) ([]user.Account, error) {
	accounts := make([]user.Account, 0)

	err := env.dbs.Read.Select(&accounts,
		user.StmtListAccounts,
		p.Limit,
		p.Offset())

	if err != nil {
		return accounts, err
	}

	return accounts, nil
}

func (env Env) ListStaff(p gorest.Pagination) (pkg.PagedList[user.Account], error) {
	countCh := make(chan int64)
	listCh := make(chan pkg.AsyncResult[[]user.Account])

	go func() {
		defer close(countCh)

		n, err := env.countStaff()
		if err != nil {
			log.Print(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)

		list, err := env.listStaff(p)
		if err != nil {
			log.Print(err)
		}

		listCh <- pkg.AsyncResult[[]user.Account]{
			Err:   err,
			Value: list,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return pkg.PagedList[user.Account]{}, listResult.Err
	}

	return pkg.PagedList[user.Account]{
		Total:      count,
		Pagination: p,
		Data:       listResult.Value,
	}, nil
}
