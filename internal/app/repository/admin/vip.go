package admin

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg"
	"github.com/FTChinese/superyard/pkg/reader"
	"log"
)

// FtcAccount retrieves an ftc account before granting/revoking vip.
func (env Env) FtcAccount(ftcID string) (reader.BaseAccount, error) {
	var a reader.BaseAccount
	err := env.dbs.Read.Get(
		&a,
		reader.StmtAccountByFtc,
		ftcID)
	if err != nil {
		return a, err
	}

	return a, nil
}

func (env Env) countVip() (int64, error) {
	var count int64
	err := env.dbs.Read.Get(&count, reader.StmtCountVIP)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listVip(p gorest.Pagination) ([]reader.BaseAccount, error) {
	var vips = make([]reader.BaseAccount, 0)
	err := env.dbs.Read.Select(&vips, reader.StmtListVIP, p.Limit, p.Offset())
	if err != nil {
		return nil, err
	}

	return vips, nil
}

func (env Env) ListVIP(p gorest.Pagination) (pkg.PagedList[reader.BaseAccount], error) {
	countCh := make(chan int64)
	listCh := make(chan pkg.AsyncResult[[]reader.BaseAccount])

	go func() {
		defer close(countCh)
		n, err := env.countVip()

		if err != nil {
			log.Print(err)
		}

		countCh <- n
	}()

	go func() {
		defer close(listCh)
		list, err := env.listVip(p)
		listCh <- pkg.AsyncResult[[]reader.BaseAccount]{
			Err:   err,
			Value: list,
		}
	}()

	count, listResult := <-countCh, <-listCh

	if listResult.Err != nil {
		return pkg.PagedList[reader.BaseAccount]{}, listResult.Err
	}

	return pkg.PagedList[reader.BaseAccount]{
		Total:      count,
		Pagination: p,
		Data:       listResult.Value,
	}, nil
}

// UpdateVIP set/removes vip column.
func (env Env) UpdateVIP(a reader.BaseAccount) error {
	_, err := env.dbs.Read.NamedExec(reader.StmtSetVIP, a)

	if err != nil {
		return err
	}

	return nil
}
