package admin

import (
	gorest "github.com/FTChinese/go-rest"
	"github.com/FTChinese/superyard/pkg/reader"
	"log"
)

// FtcAccount retrieves an ftc account before granting/revoking vip.
func (env Env) FtcAccount(ftcID string) (reader.FtcAccount, error) {
	var a reader.FtcAccount
	err := env.db.Get(&a, reader.StmtFtcAccount, ftcID)
	if err != nil {
		return a, err
	}

	return a, nil
}

func (env Env) countVip() (int64, error) {
	var count int64
	err := env.db.Get(&count, reader.StmtCountVIP)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (env Env) listVip(p gorest.Pagination) ([]reader.FtcAccount, error) {
	var vips = make([]reader.FtcAccount, 0)
	err := env.db.Select(&vips, reader.StmtListVIP, p.Limit, p.Offset())
	if err != nil {
		return nil, err
	}

	return vips, nil
}

func (env Env) ListVIP(p gorest.Pagination) (reader.FtcAccountList, error) {
	countCh := make(chan int64)
	listCh := make(chan reader.FtcAccountList)

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

// UpdateVIP set/removes vip column.
func (env Env) UpdateVIP(a reader.FtcAccount) error {
	_, err := env.db.NamedExec(reader.StmtSetVIP, a)

	if err != nil {
		return err
	}

	return nil
}
